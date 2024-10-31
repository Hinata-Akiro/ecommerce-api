package routes

import (
	"ecommerce-api/auth"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuthSetUpRoute sets up authentication routes for registration and login
func AuthSetUpRoute(router *gin.RouterGroup, jwtSecret string, db *gorm.DB) {
	authService := auth.NewAuthService(jwtSecret, db)
	authController := auth.NewAuthController(authService)

	auth := router.Group("/auth")

	auth.POST("/register", authController.Register)
	auth.POST("/login", authController.Login)
}
