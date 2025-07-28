package api

import (
	"database/sql"
	"net/http"
	"strconv"

	db "github.com/Stylozone/go-ecom-api/db/sqlc"
	"github.com/Stylozone/go-ecom-api/dto"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	Store db.Querier
}

func NewProductHandler(store db.Querier) *ProductHandler {
	return &ProductHandler{Store: store}
}

func (h *ProductHandler) RegisterRoutes(r *gin.Engine, auth gin.HandlerFunc, admin gin.HandlerFunc) {
	public := r.Group("/products")
	public.GET("/", h.ListProducts)
	public.GET("/:id", h.GetProduct)

	protected := r.Group("/products")
	protected.Use(auth, admin)
	protected.POST("/", h.CreateProduct)
	// protected.PUT("/:id", h.UpdateProduct)
	// protected.DELETE("/:id", h.DeleteProduct)
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req dto.CreateProductRequest
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
	total, _ := h.Store.CountProducts(c)
	c.JSON(http.StatusOK, gin.H{
		"items": products,
		"total": total,
	})
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
