package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
)

var (
	ErrUnAuthorized       = errors.New("Unauthorized")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserExists         = errors.New("user with email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Map = map[string]interface{}

// ErrorResponse is used to structure error responses
type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    *Map   `json:"data,omitempty"`
}

// APIResponse is used to structure standard API responses
type APIResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// NewAPIResponse generates a structured API response
func NewAPIResponse(status int, message string, data interface{}, err string) *APIResponse {
	return &APIResponse{
		Status:  status,
		Message: message,
		Data:    data,
		Error:   err,
	}
}

// NewErrorResponse creates an error response in a structured format
func NewErrorResponse(status int, err error) *ErrorResponse {
	return &ErrorResponse{
		Status:  status,
		Message: err.Error(),
		Data:    nil,
	}
}

// Send sends a success APIResponse to the Gin context
func (resp *APIResponse) Send(c *gin.Context) {
	c.JSON(resp.Status, resp)
}

// SendError sends an ErrorResponse to the Gin context
func (resp *ErrorResponse) SendError(c *gin.Context) {
	c.JSON(resp.Status, resp)
}
