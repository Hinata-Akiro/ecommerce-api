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

// CreateProduct creates a new product in the database.
//
// The function takes a single parameter:
// - productDTO: A pointer to a CreateProduct struct representing the product data to be created.
//   The CreateProduct struct should contain the Name, Description, Price, and Stock fields.
//
// The function creates a new Product struct using the provided data and inserts it into the database.
// It returns an error if any issues occur during the creation process.
// If the product is successfully created, the function returns nil.
func (s *ProductService) CreateProduct(productDTO *CreateProduct) error {
	product := models.Product{
		Name:        productDTO.Name,
		Description: productDTO.Description,
		Price:       productDTO.Price,
		Stock:       productDTO.Stock,
	}

	return s.db.Create(&product).Error
}


// GetProduct retrieves a product by ID from the database.
//
// The function takes a single parameter:
// - id: A string representing the unique identifier of the product to be retrieved.
//
// The function returns two values:
// - A pointer to a models.Product struct representing the retrieved product.
//   If the product is not found, the function returns nil.
// - An error, which is nil if the product is successfully retrieved.
//   If the product is not found, the error message will be "product not found".
//   If there is an error while interacting with the database, the error message will start with "database error:".
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


// ListProducts retrieves all products from the database.
//
// The function selects only the id, name, price, description, and stock fields from the products table.
// It returns a slice of Product structs and an error if any issues occur during the retrieval process.
//
// If the retrieval is successful, the function returns a slice of Product structs and nil as the error.
// If there is an error while interacting with the database, the function returns nil as the slice of Product structs
// and an error with a descriptive message.
func (s *ProductService) ListProducts() ([]models.Product, error) {
    var products []models.Product
    if err := s.db.Select("id", "name", "price", "description", "stock").
        Order("created_at DESC").
        Find(&products).Error; err != nil {
        return nil, errors.New("failed to retrieve products: " + err.Error())
    }
    return products, nil
}


// UpdateProduct updates an existing product in the database.
//
// The function takes a single parameter:
// - product: A pointer to a models.Product struct representing the product to be updated.
//
// The function first checks if a product with the given ID exists in the database.
// If the product is not found, it returns an error with the message "product not found".
// If the product is found, it updates the existing product record with the data from the provided product struct.
//
// The function returns an error if any issues occur during the update process.
// If the product is successfully updated, the function returns nil.
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


// DeleteProduct deletes a product by ID from the database.
//
// The function takes a single parameter:
// - id: A string representing the unique identifier of the product to be deleted.
//
// The function returns an error if any issues occur during the deletion process.
// If the product is successfully deleted, the function returns nil.
// If the product with the given ID does not exist, the function returns an error with the message "product not found".
// If there is an error while interacting with the database, the function returns an error with a descriptive message.
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

