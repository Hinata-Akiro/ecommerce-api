package products

import (
	"ecommerce-api/models"
	"errors"

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
func (s *ProductService) CreateProduct(productDTO *CreateProduct) error {
	product := models.Product{
		Name:        productDTO.Name,
		Description: productDTO.Description,
		Price:       productDTO.Price,
		Stock:       productDTO.Stock,
	}

	return s.db.Create(&product).Error
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
	if err := s.db.Select("id", "name", "price", "description", "stock").Find(&products).Error; err != nil {
		return nil, errors.New("failed to retrieve products: " + err.Error())
	}
	return products, nil
}

// UpdateProduct updates an existing product
func (s *ProductService) UpdateProduct(product *models.Product) error {
	var existingProduct models.Product
	if err := s.db.First(&existingProduct, product.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("product not found")
		}
		return err
	}

	return s.db.Model(&existingProduct).Updates(product).Error
}

// DeleteProduct deletes a product by ID
func (s *ProductService) DeleteProduct(id string) error {
	result := s.db.Unscoped().Delete(&models.Product{}, "id = ?", id)
	if result.Error != nil {
		return errors.New("failed to delete product: " + result.Error.Error())
	}

	if result.RowsAffected == 0 {
		return errors.New("product not found")
	}
	return nil
}
