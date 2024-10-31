package routes

import (
	"ecommerce-api/middleware"
	"ecommerce-api/orders"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// OrderSetUpRoute sets up routes for order management
func OrderSetUpRoute(router *gin.RouterGroup, db *gorm.DB) {
	orderService := orders.NewOrderService(db)
	orderController := orders.NewOrderController(orderService)

	order := router.Group("/orders")
	order.Use(middleware.AuthMiddleware()) 

	order.POST("", orderController.PlaceOrder)
	order.GET("", orderController.ListOrders)
	order.PUT("/:id/cancel", orderController.CancelOrder)

	order.PUT("/:id/status", middleware.AdminMiddleware(), orderController.UpdateOrderStatus)
}
