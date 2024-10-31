package products

type CreateProduct struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Price       int64  `json:"price" binding:"required"`
	Stock       int    `json:"stock" binding:"required"`
}

type UpdateProduct struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Price       *int64  `json:"price,omitempty"`
	Stock       *int    `json:"stock,omitempty"`
}
