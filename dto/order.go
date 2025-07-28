package dto

type CreateOrderItem struct {
	ProductID int32 `json:"product_id" binding:"required"`
	Quantity  int32 `json:"quantity" binding:"required"`
}

type CreateOrderRequest struct {
	Items []CreateOrderItem `json:"items" binding:"required,dive"`
}

type OrderItemResponse struct {
	ProductID int32 `json:"product_id"`
	Quantity  int32 `json:"quantity"`
	Price     int32 `json:"price"`
}

type OrderResponse struct {
	OrderID    int32               `json:"order_id"`
	TotalPrice int32               `json:"total_price"`
	CreatedAt  string              `json:"created_at"`
	Items      []OrderItemResponse `json:"items"`
}
