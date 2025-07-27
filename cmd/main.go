package main

import (
	"log"

	"github.com/Stylozone/go-ecom-api/pkg/config"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := config.LoadConfig("."); err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"env":    config.AppConfig.Port,
		})
	})

	log.Printf("Server running at http://localhost:%s\n", config.AppConfig.Port)
	r.Run(":" + config.AppConfig.Port)
}
