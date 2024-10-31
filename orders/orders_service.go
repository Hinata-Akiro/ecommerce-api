package orders

import (
	"ecommerce-api/models"
	"errors"

	"gorm.io/gorm"
)

type OrderService struct {
	db *gorm.DB
}

func NewOrderService(db *gorm.DB) *OrderService { return &OrderService{db: db} }

func (s *OrderService) PlaceOrder(userID uint, products []models.OrderProduct) (*models.Order, error) {
	order := models.Order{UserID: userID, Status: models.OrderStatusPending}

	if err := s.db.Create(&order).Error; err != nil {
		return nil, errors.New("failed to create order: " + err.Error())
	}

	orderProducts := make([]models.OrderProduct, len(products))
	for i, item := range products {
		orderProducts[i] = models.OrderProduct{
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
	}

	if err := s.db.CreateInBatches(&orderProducts, 100).Error; err != nil {
		return nil, errors.New("failed to create order-product associations")
	}

	if err := s.db.Preload("Products").First(&order, order.ID).Error; err != nil {
		return nil, errors.New("failed to retrieve order with products")
	}

	return &order, nil
}

// ListOrders retrieves all orders for a specific user
func (s *OrderService) ListOrders(userID uint) ([]models.Order, error) {
	var orders []models.Order
	err := s.db.Preload("Products").Where("user_id = ?", userID).Find(&orders).Error
	if err != nil {
		return nil, errors.New("failed to retrieve orders: " + err.Error())
	}
	if len(orders) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return orders, nil
}

// CancelOrder cancels a pending order if it belongs to the user
func (s *OrderService) CancelOrder(orderID, userID uint) error {
	err := s.db.Model(&models.Order{}).
		Where("id = ? AND user_id = ? AND status = ?", orderID, userID, models.OrderStatusPending).
		Update("status", models.OrderStatusCancelled).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("order not found or not eligible for cancellation")
		}
		return errors.New("failed to cancel order: " + err.Error())
	}
	return nil
}

// UpdateOrderStatus updates the status of an order (admin only)
func (s *OrderService) UpdateOrderStatus(orderID uint, status models.OrderStatus) (*models.Order, error) {
	var order models.Order
	if err := s.db.First(&order, orderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, errors.New("failed to retrieve order: " + err.Error())
	}

	order.Status = status
	if err := s.db.Save(&order).Error; err != nil {
		return nil, errors.New("failed to update order status")
	}

	if err := s.db.Preload("Products").First(&order, orderID).Error; err != nil {
		return nil, errors.New("failed to retrieve updated order with products")
	}
	return &order, nil
}
