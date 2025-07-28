package router

import (
	"github.com/Stylozone/go-ecom-api/api"
	"github.com/Stylozone/go-ecom-api/db/sqlc"
	"github.com/Stylozone/go-ecom-api/middleware"
	"github.com/Stylozone/go-ecom-api/pkg/config"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, store sqlc.Querier, cfg config.Config) {
	// Middleware
	authMiddleware := middleware.AuthMiddleware(cfg.JWTSecret)
	adminOnly := middleware.AdminOnly()

	// Health
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "env": cfg.Port})
	})

	// Auth
	authHandler := api.NewAuthHandler(store, cfg.JWTSecret)
	authHandler.RegisterRoutes(r)

	// Products
	productHandler := api.NewProductHandler(store)
	productHandler.RegisterRoutes(r, authMiddleware, adminOnly)

	// Orders
	orderHandler := api.NewOrderHandler(store)
	orderHandler.RegisterRoutes(r, authMiddleware)
}
