package main

import (
	"database/sql"
	"log"

	"github.com/Stylozone/go-ecom-api/db/sqlc"
	"github.com/Stylozone/go-ecom-api/pkg/config"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	// Load config from app.env or environment
	if err := config.LoadConfig("."); err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	// Connect to PostgreSQL
	db, err := sql.Open("postgres", config.AppConfig.DBSource)
	if err != nil {
		log.Fatalf("cannot connect to database: %v", err)
	}

	// Ping to verify connection
	if err := db.Ping(); err != nil {
		log.Fatalf("database not reachable: %v", err)
	}

	// Initialize sqlc Queries
	store := sqlc.New(db)

	// Setup Gin
	r := gin.Default()

	// Healthcheck
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"env":    config.AppConfig.Port,
		})
	})

	log.Printf("Server running at http://localhost:%s\n", config.AppConfig.Port)
	err = r.Run(":" + config.AppConfig.Port)
	if err != nil {
		log.Fatal("server failed to start:", err)
	}

	// store is ready to be injected into handlers
	_ = store
}
