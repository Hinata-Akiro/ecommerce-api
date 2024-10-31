package middleware

import (
	"ecommerce-api/config"
	"ecommerce-api/database"
	"ecommerce-api/models"
	"ecommerce-api/utils"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			utils.NewErrorResponse(http.StatusUnauthorized, errors.New("Unauthorized")).SendError(c)
			c.Abort()
			return
		}

		jwtKey := []byte(config.CONFIG.JWT_SECRET)

		token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			utils.NewErrorResponse(http.StatusUnauthorized, err).SendError(c)
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*jwt.RegisteredClaims)
		if !ok || claims.Subject == "" || claims.ExpiresAt == nil || claims.ExpiresAt.Time.Before(time.Now()) {
			utils.NewErrorResponse(http.StatusUnauthorized, errors.New("Unauthorized")).SendError(c)
			c.Abort()
			return
		}

		c.Set("userID", claims.Subject)
		c.Next()
	}
}


func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			utils.NewAPIResponse(http.StatusUnauthorized, "Unauthorized", nil, "User ID not found in context").Send(c)
			c.Abort()
			return
		}

		id, err := strconv.ParseUint(userID.(string), 10, 32)
		if err != nil {
			utils.NewAPIResponse(http.StatusBadRequest, "Invalid User ID", nil, "User ID conversion failed").Send(c)
			c.Abort()
			return
		}

		var user models.User
		if err := database.Database.Take(&user, "id = ? AND is_admin = ?", id, true).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.NewAPIResponse(http.StatusForbidden, "Admin privileges required", nil, "User does not have admin rights").Send(c)
			} else {
				utils.NewAPIResponse(http.StatusInternalServerError, "Failed to retrieve user", nil, err.Error()).Send(c)
			}
			c.Abort()
			return
		}

		c.Next()
	}
}
