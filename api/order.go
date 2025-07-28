package api

import (
	"net/http"
	"strconv"

	db "github.com/Stylozone/go-ecom-api/db/sqlc"
	"github.com/Stylozone/go-ecom-api/dto"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	Store db.Querier
}

func NewOrderHandler(store db.Querier) *OrderHandler {
	return &OrderHandler{Store: store}
}

func (h *OrderHandler) RegisterRoutes(group *gin.RouterGroup) {
	group.POST("/", h.CreateOrder)
	group.GET("/user/:id", h.GetOrdersByUser)
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in token"})
		return
	}

	total := int32(0)
	orderItems := make([]db.CreateOrderItemParams, 0, len(req.Items))

	for _, item := range req.Items {
		product, err := h.Store.GetProduct(c, item.ProductID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
			return
		}

		subtotal := product.Price * item.Quantity
		total += subtotal

		orderItems = append(orderItems, db.CreateOrderItemParams{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price, // use DB price
		})
	}

	order, err := h.Store.CreateOrder(c, db.CreateOrderParams{
		UserID:     userID.(int32),
		TotalPrice: total,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Add order_id to each item and insert
	for _, item := range orderItems {
		item.OrderID = order.ID
		err := h.Store.CreateOrderItem(c, item)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusCreated, order)
}

func (h *OrderHandler) GetOrdersByUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	rawOrders, err := h.Store.GetOrdersByUser(c, int32(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	orderMap := make(map[int32]*dto.OrderResponse)

	for _, row := range rawOrders {
		order, exists := orderMap[row.OrderID]
		if !exists {
			order = &dto.OrderResponse{
				OrderID:    row.OrderID,
				TotalPrice: row.TotalPrice,
				CreatedAt:  row.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
				Items:      []dto.OrderItemResponse{},
			}
			orderMap[row.OrderID] = order
		}

		order.Items = append(order.Items, dto.OrderItemResponse{
			ProductID: row.ProductID,
			Quantity:  row.Quantity,
			Price:     row.Price,
		})
	}

	// Flatten map to slice
	result := make([]dto.OrderResponse, 0, len(orderMap))
	for _, o := range orderMap {
		result = append(result, *o)
	}

	c.JSON(http.StatusOK, result)
}
