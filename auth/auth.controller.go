package auth

import (
	"ecommerce-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

type AuthController struct {
	authService *AuthService
}

// NewAuthController initializes a new AuthController
func NewAuthController(authService *AuthService) *AuthController {
	validate = validator.New()
	return &AuthController{authService: authService}
}

// Register godoc
// @Summary      Register a new user
// @Description  Registers a new user with email and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body      RegisterDTO   true  "User registration details"
// @Success      200    {object}  utils.APIResponse
// @Failure      400    {object}  utils.APIResponse
// @Failure      500    {object}  utils.APIResponse
// @Router       /auth/register [post]
func (c *AuthController) Register(ctx *gin.Context) {
	var user RegisterDTO
	if err := ctx.ShouldBindJSON(&user); err != nil {
		utils.NewAPIResponse(http.StatusBadRequest, "Invalid input", nil, err.Error()).Send(ctx)
		return
	}

	if err := validate.Struct(&user); err != nil {
		utils.NewAPIResponse(http.StatusBadRequest, "Validation failed", nil, err.Error()).Send(ctx)
		return
	}

	if err := c.authService.Register(&user); err != nil {
		utils.NewAPIResponse(http.StatusInternalServerError, "Failed to register user", nil, err.Error()).Send(ctx)
		return
	}

	utils.NewAPIResponse(http.StatusOK, "User registered successfully", nil, "").Send(ctx)
}

// Login godoc
// @Summary      Login a user
// @Description  Authenticates a user and returns a JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body      LoginDTO   true  "User login credentials"
// @Success      200    {object}  utils.APIResponse{data=string} "JWT Token"
// @Failure      400    {object}  utils.APIResponse
// @Failure      401    {object}  utils.APIResponse
// @Failure      404    {object}  utils.APIResponse
// @Failure      500    {object}  utils.APIResponse
// @Router       /auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var input LoginDTO

	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.NewAPIResponse(http.StatusBadRequest, "Invalid input", nil, err.Error()).Send(ctx)
		return
	}

	if err := validate.Struct(&input); err != nil {
		utils.NewAPIResponse(http.StatusBadRequest, "Validation failed", nil, err.Error()).Send(ctx)
		return
	}

	token, err := c.authService.Login(input.Email, input.Password)
	if err != nil {
		switch err {
		case utils.ErrUserNotFound:
			utils.NewAPIResponse(http.StatusNotFound, "User not found", nil, "").Send(ctx)
		case utils.ErrInvalidCredentials:
			utils.NewAPIResponse(http.StatusUnauthorized, "Invalid credentials", nil, "").Send(ctx)
		default:
			utils.NewAPIResponse(http.StatusInternalServerError, "Failed to generate token", nil, "").Send(ctx)
		}
		return
	}

	utils.NewAPIResponse(http.StatusOK, "Login successful", gin.H{"token": token}, "").Send(ctx)
}
