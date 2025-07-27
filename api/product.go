package api

import (
	"database/sql"
	"net/http"
	"strconv"

	db "github.com/Stylozone/go-ecom-api/db/sqlc"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	Store db.Querier
}

func NewProductHandler(store db.Querier) *ProductHandler {
	return &ProductHandler{Store: store}
}

func (h *ProductHandler) RegisterRoutes(r *gin.Engine) {
	group := r.Group("/products")
	group.POST("/", h.CreateProduct)
	group.GET("/", h.ListProducts)
	group.GET("/:id", h.GetProduct)
}

type createProductRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Price       int32  `json:"price" binding:"required"`
	Quantity    int32  `json:"quantity" binding:"required"`
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req createProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.Store.CreateProduct(c, db.CreateProductParams{
		Name: req.Name,
		Description: sql.NullString{
			String: req.Description,
			Valid:  req.Description != "",
		},
		Price:    req.Price,
		Quantity: req.Quantity,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

func (h *ProductHandler) ListProducts(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	products, err := h.Store.ListProducts(c, db.ListProductsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	product, err := h.Store.GetProduct(c, int32(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}
