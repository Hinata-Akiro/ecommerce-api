package routes

import (
	"ecommerce-api/middleware"
	"ecommerce-api/products"
     
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ProductSetUpRoute sets up routes for product management
func ProductSetUpRoute(router *gin.RouterGroup, db *gorm.DB) {
	productService := products.NewProductService(db)
	productController := products.NewProductController(productService)

	// Product routes group with authentication middleware
	product := router.Group("/products")
	product.Use(middleware.AuthMiddleware()) 

	// Admin-only routes
	product.POST("", middleware.AdminMiddleware(), productController.CreateProduct) // Admin-only
	product.PUT("/:id", middleware.AdminMiddleware(), productController.UpdateProduct) // Admin-only
	product.DELETE("/:id", middleware.AdminMiddleware(), productController.DeleteProduct) // Admin-only

	// Routes accessible to any authenticated user
	product.GET("/:id", productController.GetProduct)
	product.GET("", productController.ListProducts)
}
