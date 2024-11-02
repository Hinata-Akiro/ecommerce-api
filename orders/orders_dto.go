package orders

import (
	"ecommerce-api/models"
)

type ProductOrder struct {
	ProductID uint `json:"productID" binding:"required,gt=0"` 
	Quantity  int  `json:"quantity" binding:"required,gte=1"`
}

type PlaceOrderDTO struct {
	Products []ProductOrder `json:"products"`
}

type UpdateOrderStatusDTO struct {
	Status models.OrderStatus `json:"status" binding:"required,orderStatus"`
}

type OrderSummary struct {
	ID           uint    `json:"id"`
	ProductName  string  `json:"product_name"`
	Description  string  `json:"description"`
	ProductPrice float64 `json:"product_price"`
	Quantity     int     `json:"quantity"`
	TotalPrice   float64 `json:"total_price"`
}
