package router

import (
	"time"

	"github.com/Stylozone/go-ecom-api/api"
	"github.com/Stylozone/go-ecom-api/db/sqlc"
	"github.com/Stylozone/go-ecom-api/middleware"
	"github.com/Stylozone/go-ecom-api/pkg/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, store sqlc.Querier, cfg config.Config) {
	// register cors
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.FrontendURL},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

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

	// Users
	userHandler := api.NewUserHandler(store)
	userHandler.RegisterRoutes(r, authMiddleware)
}
