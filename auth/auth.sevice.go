// auth/auth_service.go
package auth

import (
	"errors"
	"time"
	"strconv"

	"ecommerce-api/models"
	"ecommerce-api/utils"

	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

// AuthService struct holds the database connection and JWT secret for auth operations
type AuthService struct {
	jwtSecret string
	db        *gorm.DB
}

// NewAuthService initializes AuthService with jwtSecret and database connection
func NewAuthService(jwtSecret string, db *gorm.DB) *AuthService {
	return &AuthService{jwtSecret: jwtSecret, db: db}
}

// Register creates a new user with a hashed password.
//
// The function takes a pointer to a User struct as a parameter. It first hashes the user's password
// using the HashPassword function from the utils package. If the hashing process fails, it returns
// the error. Otherwise, it sets the hashed password in the User struct.
//
// The function then uses a database transaction to ensure atomicity. It creates a new user record
// in the database using the provided User struct. If the creation process fails, it returns an error
// with a descriptive message. If the creation is successful, it returns nil.
func (s *AuthService) Register(userDTO *RegisterDTO) error {
	hashedPassword, err := utils.HashPassword(userDTO.Password)
	if err != nil {
		return err
	}

	user := models.User{
		Email:    userDTO.Email,
		Name:     userDTO.Name,
		Password: hashedPassword,
		IsAdmin:  false, 
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return errors.New("failed to create user: " + err.Error())
		}
		return nil
	})
}



// Login attempts to authenticate a user with the provided email and password.
// It retrieves the user from the database using the provided email.
// If the user is found, it checks if the provided password matches the stored hashed password.
// If the credentials are valid, it generates a JWT token using the user's ID as the subject.
// The token is valid for 24 hours.
//
// Parameters:
// - email: The email of the user attempting to authenticate.
// - password: The password provided by the user.
//
// Returns:
// - A string representing the JWT token if the authentication is successful.
// - An error if the authentication fails or encounters any database or token generation errors.
func (s *AuthService) Login(email, password string) (string, error) {
	var user models.User
	// Find user by email
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", utils.ErrUserNotFound
		}
		return "", errors.New("database error: " + err.Error())
	}

	// Check if the provided password matches the stored hashed password
	if err := utils.ComparePasswords(password, user.Password); err != nil {
		return "", utils.ErrInvalidCredentials
	}

	// Generate JWT token using the user's ID as the subject
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   strconv.FormatUint(uint64(user.ID), 10),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	})
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return tokenString, nil
}

