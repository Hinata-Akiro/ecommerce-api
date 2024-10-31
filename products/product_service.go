package products

import (
	"errors"
	"ecommerce-api/models"
	"gorm.io/gorm"
)

// ProductService struct manages product business logic
type ProductService struct {
	db *gorm.DB
}

// NewProductService initializes ProductService with database connection
func NewProductService(db *gorm.DB) *ProductService {
	return &ProductService{db: db}
}

// CreateProduct creates a new product in the database
func (s *ProductService) CreateProduct(product *CreateProduct) error {
	return s.db.Omit("CreatedAt", "UpdatedAt").Create(product).Error
}

// GetProduct retrieves a product by ID
func (s *ProductService) GetProduct(id string) (*models.Product, error) {
	var product models.Product
	if err := s.db.First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, errors.New("database error: " + err.Error())
	}
	return &product, nil
}

// ListProducts retrieves all products
func (s *ProductService) ListProducts() ([]models.Product, error) {
	var products []models.Product
	if err := s.db.Select("id", "name", "price").Find(&products).Error; err != nil {
		return nil, errors.New("failed to retrieve products: " + err.Error())
	}
	return products, nil
}

// UpdateProduct updates an existing product
func (s *ProductService) UpdateProduct(product *models.Product) error {
	return s.db.Model(product).Updates(product).Error
}

// DeleteProduct deletes a product by ID
func (s *ProductService) DeleteProduct(id string) error {
	if err := s.db.Unscoped().Delete(&models.Product{}, "id = ?", id).Error; err != nil {
		return errors.New("failed to delete product: " + err.Error())
	}
	return nil
}
