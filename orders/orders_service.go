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

// PlaceOrder creates a new order for the specified user and products.
//
// userID: The unique identifier of the user placing the order.
// products: A slice of ProductOrder structs representing the products to be included in the order.
//
// The function returns a pointer to the created Order struct and an error if any occurred during the process.
// If the products in the order are not found or if there's an issue with the database, an error will be returned.
//
// The function performs the following steps:
// 1. Extracts the product IDs from the input products slice.
// 2. Validates the products by querying the database for the specified product IDs.
// 3. Creates a new Order struct with the provided user ID and a pending status.
// 4. Inserts the new order into the database.
// 5. Creates OrderProduct associations for each product in the order.
// 6. Retrieves the created order with its associated products from the database.
//
// If any error occurs during the process, the function returns nil for the Order pointer and an error describing the issue.
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


// ListOrders retrieves all orders for a specific user.
//
// userID: The unique identifier of the user whose orders are to be retrieved.
//
// The function returns a slice of OrderSummary structs and an error if any occurred during the process.
// If no orders are found for the specified user, the function returns nil for the slice and gorm.ErrRecordNotFound.
// If there's an issue with the database, an error will be returned.
//
// The function performs the following steps:
// 1. Initializes an empty slice of OrderSummary structs.
// 2. Executes a database query to retrieve order details, including product information and total prices.
// 3. Checks for any errors during the query execution.
// 4. If no errors occur, checks if any orders were found for the specified user.
// 5. Returns the slice of OrderSummary structs and nil for the error if orders were found.
// 6. Returns nil for the slice and gorm.ErrRecordNotFound if no orders were found.
// 7. Returns nil for the slice and an error if any other database error occurs.
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


// CancelOrder cancels a pending order if it belongs to the user.
//
// Parameters:
// - orderID: The unique identifier of the order to be cancelled.
// - userID: The unique identifier of the user who is attempting to cancel the order.
//
// Return:
// - An error if any occurred during the cancellation process.
//   - Returns nil if the order was successfully cancelled.
//   - Returns an error with the message "order not found" if the order could not be found.
//   - Returns an error with the message "failed to retrieve order: <error details>" if there was an issue retrieving the order from the database.
//   - Returns an error with the message "order is not eligible for cancellation" if the order is not in a pending status.
//   - Returns an error with the message "failed to cancel order: <error details>" if there was an issue updating the order status in the database.
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


// UpdateOrderStatus updates the status of an order (admin only).
//
// Parameters:
// - orderID: The unique identifier of the order to be updated.
// - status: The new status to be set for the order.
//
// Return:
// - A pointer to the updated Order struct if the operation is successful.
// - An error if any occurred during the update process.
//   - Returns nil if the order was successfully updated.
//   - Returns an error with the message "order not found" if the order could not be found.
//   - Returns an error with the message "failed to retrieve order: <error details>" if there was an issue retrieving the order from the database.
//   - Returns an error with the message "failed to update order status" if there was an issue updating the order status in the database.
//   - Returns an error with the message "failed to retrieve updated order with products" if there was an issue retrieving the updated order with its associated products.
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

