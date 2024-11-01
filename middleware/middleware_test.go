package middleware

import (
	"ecommerce-api/config"
	"ecommerce-api/database"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Mock JWT Secret for testing
func init() {
	if config.CONFIG == nil {
		config.CONFIG = &config.AppConfig{}
	}
	config.CONFIG.JWT_SECRET = "testsecret"
}

func generateTestJWT(userID string, expiration time.Duration) (string, error) {
	if config.CONFIG == nil || config.CONFIG.JWT_SECRET == "" {
		return "", errors.New("empty jwt secret")
	}

	claims := &jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.CONFIG.JWT_SECRET))
}


// Test AuthMiddleware with valid and invalid tokens
func TestAuthMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		token          string
		expectedStatus int
	}{
		{
			name:           "Valid Token",
			token:          func() string { token, _ := generateTestJWT("1", time.Hour); return token }(),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Expired Token",
			token:          func() string { token, _ := generateTestJWT("1", -time.Hour); return token }(),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Missing Token",
			token:          "",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.Default()
			router.Use(AuthMiddleware())
			router.GET("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "Success"})
			})

			req, _ := http.NewRequest(http.MethodGet, "/test", nil)
			if tt.token != "" {
				req.Header.Set("Authorization", tt.token)
			}
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.NotNil(t, w)
			assert.NotNil(t, req)
			assert.NotNil(t, router)
			assert.NotNil(t, tt)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}


func TestAdminMiddleware(t *testing.T) {
	// Set up sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	// Initialize GORM with sqlmock database
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to initialize gorm with sqlmock: %s", err)
	}

	// Replace the global Database variable with our mock database
	database.Database = gormDB

	// Test cases
	tests := []struct {
		name           string
		userID         string
		isAdmin        bool
		expectedStatus int
		mockError      error
	}{
		{"User ID not found in context", "", false, http.StatusUnauthorized, nil},
		{"Invalid User ID", "invalid-id", false, http.StatusBadRequest, nil},
		{"User not found in database", "123", false, http.StatusForbidden, gorm.ErrRecordNotFound},
		{"User without admin rights", "123", false, http.StatusForbidden, nil},
		{"Admin User", "123", true, http.StatusOK, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up the test request and context
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest("GET", "/", nil)
			if tt.userID != "" {
				c.Set("userID", tt.userID)
			}

			// Mock database responses based on the test case
			userIDInt, err := strconv.Atoi(tt.userID)
			if tt.userID != "" && err == nil {
				if tt.expectedStatus == http.StatusOK && tt.isAdmin {
					// Admin user found
					mock.ExpectQuery(`SELECT * FROM "users" WHERE id = $1 AND is_admin = $2`).
						WithArgs(userIDInt, true).
						WillReturnRows(sqlmock.NewRows([]string{"id", "is_admin"}).AddRow(userIDInt, true))
				} else if tt.expectedStatus == http.StatusForbidden && !tt.isAdmin {
					// Non-admin user found
					mock.ExpectQuery(`SELECT * FROM "users" WHERE id = $1 AND is_admin = $2`).
						WithArgs(userIDInt, true).
						WillReturnRows(sqlmock.NewRows([]string{"id", "is_admin"}).AddRow(userIDInt, false))
				} else if tt.mockError == gorm.ErrRecordNotFound {
					// User not found
					mock.ExpectQuery(`SELECT * FROM "users" WHERE id = $1 AND is_admin = $2`).
						WithArgs(userIDInt, true).
						WillReturnError(gorm.ErrRecordNotFound)
				}
			}

			// Call the middleware
			AdminMiddleware()(c)

			// Assert the expected HTTP status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Ensure all expectations were met
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled sqlmock expectations: %s", err)
			}
		})
	}
}

