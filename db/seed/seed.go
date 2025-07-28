package seed

import (
	"context"
	"database/sql"
	"log"
	"math/rand"
	"time"

	db "github.com/Stylozone/go-ecom-api/db/sqlc"
	"github.com/brianvoe/gofakeit/v6"
)

func SeedProducts(store *db.Queries) {
	gofakeit.Seed(time.Now().UnixNano())

	for i := 0; i < 20; i++ {
		name := gofakeit.ProductName()
		description := sql.NullString{String: gofakeit.Sentence(8), Valid: true}
		price := int32(gofakeit.Price(500, 50000))
		quantity := int32(rand.Intn(50) + 1)
		image := sql.NullString{String: gofakeit.ImageURL(640, 480), Valid: true}

		params := db.CreateProductParams{
			Name:        name,
			Description: description,
			Price:       price,
			Quantity:    quantity,
			ImageUrl:    image,
		}

		_, err := store.CreateProduct(context.Background(), params)
		if err != nil {
			log.Printf("❌ Failed to insert product: %v", err)
		} else {
			log.Printf("✅ Inserted: %s", name)
		}
	}
}
