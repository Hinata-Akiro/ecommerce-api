package products

import (
	"ecommerce-api/models"
	"ecommerce-api/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ProductController struct to handle HTTP requests for products
type ProductController struct {
	productService *ProductService
}

// NewProductController initializes a new ProductController
func NewProductController(productService *ProductService) *ProductController {
	return &ProductController{productService: productService}
}

// CreateProduct godoc
// @Summary      Create a new product
// @Description  Admin only: Creates a new product
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        product  body      CreateProduct   true  "Product details"
// @Success      201      {object}  utils.APIResponse{data=models.Product}
// @Failure      400      {object}  utils.APIResponse
// @Failure      500      {object}  utils.APIResponse
// @Security     BearerAuth
// @Router       /products [post]
func (c *ProductController) CreateProduct(ctx *gin.Context) {
	var product CreateProduct
	if err := ctx.ShouldBindJSON(&product); err != nil {
		utils.NewAPIResponse(http.StatusBadRequest, "Invalid input", nil, err.Error()).Send(ctx)
		return
	}

	if err := c.productService.CreateProduct(&product); err != nil {
		utils.NewAPIResponse(http.StatusInternalServerError, "Failed to create product", nil, err.Error()).Send(ctx)
		return
	}

	utils.NewAPIResponse(http.StatusCreated, "Product created successfully", product, "").Send(ctx)
}

// GetProduct godoc
// @Summary      Get a product
// @Description  Retrieve a product by ID
// @Tags         products
// @Produce      json
// @Param        id   path      string  true  "Product ID"
// @Success      200  {object}  utils.APIResponse{data=models.Product}
// @Failure      404  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Security     BearerAuth
// @Router       /products/{id} [get]
func (c *ProductController) GetProduct(ctx *gin.Context) {
	id := ctx.Param("id")
	product, err := c.productService.GetProduct(id)
	if err != nil {
		if err.Error() == "product not found" {
			utils.NewAPIResponse(http.StatusNotFound, "Product not found", nil, "").Send(ctx)
		} else {
			utils.NewAPIResponse(http.StatusInternalServerError, "Failed to retrieve product", nil, err.Error()).Send(ctx)
		}
		return
	}

	utils.NewAPIResponse(http.StatusOK, "Product retrieved successfully", product, "").Send(ctx)
}

// ListProducts godoc
// @Summary      List all products
// @Description  Retrieve all products
// @Tags         products
// @Produce      json
// @Success      200  {object}  utils.APIResponse{data=[]models.Product}
// @Failure      500  {object}  utils.APIResponse
// @Security     BearerAuth
// @Router       /products [get]
func (c *ProductController) ListProducts(ctx *gin.Context) {
	products, err := c.productService.ListProducts()
	if err != nil {
		utils.NewAPIResponse(http.StatusInternalServerError, "Failed to retrieve products", nil, err.Error()).Send(ctx)
		return
	}

	utils.NewAPIResponse(http.StatusOK, "Products retrieved successfully", products, "").Send(ctx)
}

// UpdateProduct godoc
// @Summary      Update a product
// @Description  Admin only: Update a product by ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id        path      string          true  "Product ID"
// @Param        product   body      UpdateProduct  true  "Updated product details"
// @Success      200       {object}  utils.APIResponse{data=models.Product}
// @Failure      400       {object}  utils.APIResponse
// @Failure      404       {object}  utils.APIResponse
// @Failure      500       {object}  utils.APIResponse
// @Security     BearerAuth
// @Router       /products/{id} [put]
func (c *ProductController) UpdateProduct(ctx *gin.Context) {
	idStr := ctx.Param("id")
	idUint, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.NewAPIResponse(http.StatusBadRequest, "Invalid ID", nil, err.Error()).Send(ctx)
		return
	}

	var product UpdateProduct
	if err := ctx.ShouldBindJSON(&product); err != nil {
		utils.NewAPIResponse(http.StatusBadRequest, "Invalid input", nil, err.Error()).Send(ctx)
		return
	}

	var updatedProduct models.Product
	updatedProduct.ID = uint(idUint)

	if product.Name != nil {
		updatedProduct.Name = *product.Name
	}
	if product.Description != nil {
		updatedProduct.Description = *product.Description
	}
	if product.Price != nil {
		updatedProduct.Price = *product.Price
	}
	if product.Stock != nil {
		updatedProduct.Stock = *product.Stock
	}

	err = c.productService.UpdateProduct(&updatedProduct)
	if err != nil {
		if err.Error() == "product not found" {
			utils.NewAPIResponse(http.StatusNotFound, "Product not found", nil, "").Send(ctx)
		} else {
			utils.NewAPIResponse(http.StatusInternalServerError, "Failed to update product", nil, err.Error()).Send(ctx)
		}
		return
	}

	utils.NewAPIResponse(http.StatusOK, "Product updated successfully", updatedProduct, "").Send(ctx)
}


// DeleteProduct godoc
// @Summary      Delete a product
// @Description  Admin only: Delete a product by ID
// @Tags         products
// @Produce      json
// @Param        id   path      string  true  "Product ID"
// @Success      200  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Security     BearerAuth
// @Router       /products/{id} [delete]
func (c *ProductController) DeleteProduct(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.productService.DeleteProduct(id); err != nil {
		utils.NewAPIResponse(http.StatusNotFound, "Product not found", nil, "").Send(ctx)
		return
	}

	utils.NewAPIResponse(http.StatusOK, "Product deleted successfully", nil, "").Send(ctx)
}
