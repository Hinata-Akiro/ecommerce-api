package models

import (
	"errors"
	"time"
)

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusShipped    OrderStatus = "shipped"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
)

type Order struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt *time.Time     `json:"deleted_at" gorm:"index"`
	UserID    uint           `json:"user_id"`
	Products  []Product      `gorm:"many2many:order_products;" json:"products"`
	Status    OrderStatus    `json:"status" gorm:"default:'pending'"`
}

type OrderProduct struct {
	OrderID   uint `json:"order_id"`
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}


func (status OrderStatus) IsValid() error {
	switch status {
	case OrderStatusPending, OrderStatusProcessing, OrderStatusShipped, OrderStatusDelivered, OrderStatusCancelled:
		return nil
	}
	return errors.New("invalid order status")
}