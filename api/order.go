package api

import (
	"net/http"
	"strconv"

	db "github.com/Stylozone/go-ecom-api/db/sqlc"
	"github.com/Stylozone/go-ecom-api/dto"
	"github.com/Stylozone/go-ecom-api/pkg/utils"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	Store db.Querier
}

func NewOrderHandler(store db.Querier) *OrderHandler {
	return &OrderHandler{Store: store}
}

func (h *OrderHandler) RegisterRoutes(r *gin.Engine, auth gin.HandlerFunc, admin gin.HandlerFunc) {
	protected := r.Group("/orders")
	protected.Use(auth)

	protected.POST("/", h.CreateOrder)
	protected.GET("/me", h.GetMyOrders)

	protectedAdmin := r.Group("/orders")
	protectedAdmin.Use(auth, admin)
	protectedAdmin.GET("/user/:id", h.GetOrdersByUser)
	protectedAdmin.PATCH("/:id", h.UpdateOrder)
	// protectedAdmin.GET("/", h.ListAllOrders)
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

	result := utils.GroupOrderItems(rawOrders)
	c.JSON(http.StatusOK, result)
}

func (h *OrderHandler) GetMyOrders(c *gin.Context) {
	// Get user_id from JWT middleware context
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(int32)

	rawOrders, err := h.Store.GetOrdersByUser(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot get orders"})
		return
	}

	groupedOrders := utils.GroupOrderItems(rawOrders)
	c.JSON(http.StatusOK, groupedOrders)
}

func (h *OrderHandler) UpdateOrder(c *gin.Context) {
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.Store.UpdateOrderStatus(c, db.UpdateOrderStatusParams{
		ID:     int32(orderID),
		Status: req.Status,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order status updated"})
}
