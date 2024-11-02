package products

type CreateProduct struct {
	Name        string `json:"name" binding:"required,min=1"`         
	Description string `json:"description" binding:"required,min=1"` 
	Price       int64  `json:"price" binding:"required,gt=0"`         
	Stock       int    `json:"stock" binding:"required,gt=0"`         
}

type UpdateProduct struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Price       *int64  `json:"price,omitempty"`
	Stock       *int    `json:"stock,omitempty"`
}
