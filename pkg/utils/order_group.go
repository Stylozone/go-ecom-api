package utils

import (
	"github.com/Stylozone/go-ecom-api/db/sqlc"
	"github.com/Stylozone/go-ecom-api/dto"
)

func GroupOrderItems(rows []sqlc.GetOrdersByUserRow) []dto.OrderResponse {
	orderMap := make(map[int32]*dto.OrderResponse)

	for _, row := range rows {
		order, exists := orderMap[row.OrderID]
		if !exists {
			order = &dto.OrderResponse{
				OrderID:    row.OrderID,
				TotalPrice: row.TotalPrice,
				CreatedAt:  row.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
				Items:      []dto.OrderItemResponse{},
			}
			orderMap[row.OrderID] = order
		}

		order.Items = append(order.Items, dto.OrderItemResponse{
			ProductID: row.ProductID,
			Quantity:  row.Quantity,
			Price:     row.Price,
		})
	}

	result := make([]dto.OrderResponse, 0, len(orderMap))
	for _, o := range orderMap {
		result = append(result, *o)
	}
	return result
}
