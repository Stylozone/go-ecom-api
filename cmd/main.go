package main

import (
	"database/sql"
	"flag"
	"log"

	"github.com/Stylozone/go-ecom-api/db/seed"
	"github.com/Stylozone/go-ecom-api/db/sqlc"
	"github.com/Stylozone/go-ecom-api/pkg/config"
	"github.com/Stylozone/go-ecom-api/router"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func main() {
	// Add a --seed flag
	seedFlag := flag.Bool("seed", false, "Seed fake products into the database")
	flag.Parse()

	if err := config.LoadConfig("."); err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	dbConn, err := sql.Open("postgres", config.AppConfig.DBSource)
	if err != nil {
		log.Fatalf("cannot connect to database: %v", err)
	}

	if err := goose.Up(dbConn, "db/migrations"); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	if err := dbConn.Ping(); err != nil {
		log.Fatalf("database not reachable: %v", err)
	}

	store := sqlc.New(dbConn)

	if *seedFlag {
		seed.SeedProducts(store)
		return
	}

	r := gin.Default()

	// Route registration
	router.RegisterRoutes(r, store, config.AppConfig)

	log.Printf("Server running at http://localhost:%s\n", config.AppConfig.Port)
	if err := r.Run(":" + config.AppConfig.Port); err != nil {
		log.Fatal("server failed to start:", err)
	}
}
