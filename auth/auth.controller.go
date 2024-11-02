package auth

import (
	"ecommerce-api/models"
	"ecommerce-api/utils"
	"encoding/json"
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
// @Description  Registers a new user with email and password, with optional admin privileges
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body      RegisterDTO   true  "User registration details"
// @Param        admin  query     bool          false "Set to true to register user as admin"
// @Success      200    {object}  utils.APIResponse
// @Failure      400    {object}  utils.APIResponse
// @Failure      500    {object}  utils.APIResponse
// @Router       /auth/register [post]
func (c *AuthController) Register(ctx *gin.Context) {
    var userDTO RegisterDTO

    if err := ctx.ShouldBindJSON(&userDTO); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors := make(map[string]string)
			for _, v := range validationErrors {
				field := v.Field()
				if field == "Password" {
					_, errorMessage := models.ValidatePassword(userDTO.Password)
					errors["password"] = errorMessage
				} else {
					errors[field] = v.Error() 
				}
			}
			jsonErrors, _ := json.Marshal(errors)
			utils.NewAPIResponse(http.StatusBadRequest, "Validation error", nil, string(jsonErrors)).Send(ctx)
			return
		}
		utils.NewAPIResponse(http.StatusBadRequest, "Invalid input", nil, err.Error()).Send(ctx)
		return
	}

    isAdmin := false
    if adminQuery := ctx.Query("admin"); adminQuery == "true" {
        isAdmin = true
    }

    // Register the user using the auth service
    if err := c.authService.Register(&userDTO, isAdmin); err != nil {
        if err == utils.ErrUserExists {
            utils.NewAPIResponse(http.StatusBadRequest, "User with email already exists", nil, err.Error()).Send(ctx)
        } else {
            utils.NewAPIResponse(http.StatusInternalServerError, "Failed to register user", nil, err.Error()).Send(ctx)
        }
        return
    }

    // Success response
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
