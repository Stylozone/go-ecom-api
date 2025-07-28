package api

import (
	"net/http"
	"strconv"

	db "github.com/Stylozone/go-ecom-api/db/sqlc"
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

type createOrderItem struct {
	ProductID int32 `json:"product_id" binding:"required"`
	Quantity  int32 `json:"quantity" binding:"required"`
	Price     int32 `json:"price" binding:"required"`
}

type createOrderRequest struct {
	Items []createOrderItem `json:"items" binding:"required,dive"`
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req createOrderRequest
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
	for _, item := range req.Items {
		total += item.Price * item.Quantity
	}

	order, err := h.Store.CreateOrder(c, db.CreateOrderParams{
		UserID:     userID.(int32),
		TotalPrice: total,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, item := range req.Items {
		err := h.Store.CreateOrderItem(c, db.CreateOrderItemParams{
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		})
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

	orders, err := h.Store.GetOrdersByUser(c, int32(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}
