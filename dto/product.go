package dto

type CreateProductRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Price       int32  `json:"price" binding:"required"`
	Quantity    int32  `json:"quantity" binding:"required"`
}
