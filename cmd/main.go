package main

import (
	"database/sql"
	"log"

	"github.com/Stylozone/go-ecom-api/api"
	"github.com/Stylozone/go-ecom-api/db/sqlc"
	"github.com/Stylozone/go-ecom-api/pkg/config"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func main() {
	// Load config from app.env
	if err := config.LoadConfig("."); err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	// Connect to DB
	dbConn, err := sql.Open("postgres", config.AppConfig.DBSource)
	if err != nil {
		log.Fatalf("cannot connect to database: %v", err)
	}

	// ✅ Run Goose migrations here
	if err := goose.Up(dbConn, "db/migrations"); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	if err := dbConn.Ping(); err != nil {
		log.Fatalf("database not reachable: %v", err)
	}

	store := sqlc.New(dbConn) // *Queries implements Querier

	r := gin.Default()

	// Health route
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"env":    config.AppConfig.Port,
		})
	})

	// Product routes
	productHandler := api.NewProductHandler(store)
	productHandler.RegisterRoutes(r)

	// Order routes
	orderHandler := api.NewOrderHandler(store)
	orderHandler.RegisterRoutes(r)

	log.Printf("Server running at http://localhost:%s\n", config.AppConfig.Port)
	err = r.Run(":" + config.AppConfig.Port)
	if err != nil {
		log.Fatal("server failed to start:", err)
	}

}
