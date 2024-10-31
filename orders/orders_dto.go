package orders

import (
	"ecommerce-api/models"
)

type PlaceOrderDTO struct {
    Products []models.OrderProduct `json:"products"`
}

type UpdateOrderStatusDTO struct {
    Status models.OrderStatus `json:"status" binding:"required"`
}