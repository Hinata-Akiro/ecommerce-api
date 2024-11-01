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

func (s *OrderService) PlaceOrder(userID uint, products []ProductOrder) (*models.Order, error) {
	productIDs := make([]uint, len(products))
	for i, item := range products {
		productIDs[i] = item.ProductID
	}

	var dbProducts []models.Product
	if err := s.db.Where("id IN ?", productIDs).Find(&dbProducts).Error; err != nil {
		return nil, errors.New("failed to validate products: " + err.Error())
	}

	if len(dbProducts) != len(products) {
		return nil, gorm.ErrRecordNotFound
	}

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

	if err := s.db.CreateInBatches(&orderProducts, len(orderProducts)).Error; err != nil {
		return nil, errors.New("failed to create order-product associations")
	}

	if err := s.db.Preload("Products").First(&order, order.ID).Error; err != nil {
		return nil, errors.New("failed to retrieve order with products")
	}

	return &order, nil
}

// ListOrders retrieves all orders for a specific user
func (s *OrderService) ListOrders(userID uint) ([]OrderSummary, error) {
	var orderSummaries []OrderSummary
	err := s.db.
		Model(&models.Order{}).
		Select("orders.id as id, products.name as product_name, products.description as description, products.price as product_price, order_products.quantity as quantity, sum(order_products.quantity * products.price) as total_price").
		Joins("JOIN order_products ON orders.id = order_products.order_id").
		Joins("JOIN products ON products.id = order_products.product_id").
		Where("orders.user_id = ?", userID).
		Group("orders.id, products.name, products.description, products.price, order_products.quantity").
		Find(&orderSummaries).Error
	if err != nil {
		return nil, errors.New("failed to retrieve orders: " + err.Error())
	}
	if len(orderSummaries) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return orderSummaries, nil
}

// CancelOrder cancels a pending order if it belongs to the user
func (s *OrderService) CancelOrder(orderID, userID uint) error {
	var order models.Order

	if err := s.db.Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("order not found")
		}
		return errors.New("failed to retrieve order: " + err.Error())
	}

	if order.Status != models.OrderStatusPending {
		return errors.New("order is not eligible for cancellation")
	}

	if err := s.db.Model(&order).Update("status", models.OrderStatusCancelled).Error; err != nil {
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
