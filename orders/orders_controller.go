package orders

import (
	"ecommerce-api/utils"
	"ecommerce-api/models"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	 "github.com/go-playground/validator/v10"
)

var validate *validator.Validate


type OrderController struct {
	orderService *OrderService
}

func NewOrderController(orderService *OrderService) *OrderController {
	validate = validator.New()
	return &OrderController{orderService: orderService}
}

// PlaceOrder godoc
// @Summary      Place an order
// @Description  Allows a user to place an order for one or more products
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        products  body      PlaceOrderDTO   true  "List of products to order"
// @Success      201       {object}  utils.APIResponse{data=models.Order}
// @Failure      400       {object}  utils.APIResponse
// @Failure      404       {object}  utils.APIResponse
// @Failure      500       {object}  utils.APIResponse
// @Security     BearerAuth
// @Router       /orders [post]
func (c *OrderController) PlaceOrder(ctx *gin.Context) {
	var input PlaceOrderDTO
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.NewAPIResponse(http.StatusBadRequest, "Invalid input", nil, err.Error()).Send(ctx)
		return
	}

	userIDStr, exists := ctx.Get("userID")
	if !exists {
		utils.NewAPIResponse(http.StatusUnauthorized, "Unauthorized", nil, "User ID not found in context").Send(ctx)
	}
	userID, err := strconv.ParseUint(userIDStr.(string), 10, 32)
	if err != nil {
		utils.NewAPIResponse(http.StatusBadRequest, "Invalid User ID", nil, "User ID conversion failed").Send(ctx)
		return
	}
	order, err := c.orderService.PlaceOrder(uint(userID), input.Products)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			utils.NewAPIResponse(http.StatusNotFound, "One or more products do not exist", nil, "").Send(ctx)
		default:
			utils.NewAPIResponse(http.StatusInternalServerError, "Failed to place order", nil, err.Error()).Send(ctx)
		}
		return
	}

	utils.NewAPIResponse(http.StatusCreated, "Order placed successfully", order, "").Send(ctx)
}

// ListOrders godoc
// @Summary      List user's orders
// @Description  Allows a user to view their orders
// @Tags         orders
// @Produce      json
// @Success      200  {object}  utils.APIResponse{data=[]models.Order}
// @Failure      500  {object}  utils.APIResponse
// @Security     BearerAuth
// @Router       /orders [get]
func (c *OrderController) ListOrders(ctx *gin.Context) {
	// Get user ID from context
	userIDStr, exists := ctx.Get("userID")
	if !exists {
		utils.NewAPIResponse(http.StatusUnauthorized, "Unauthorized", nil, "User ID not found in context").Send(ctx)
	}
	userID, err := strconv.ParseUint(userIDStr.(string), 10, 32)
	if err != nil {
		utils.NewAPIResponse(http.StatusBadRequest, "Invalid User ID", nil, "User ID conversion failed").Send(ctx)
		return
	}

	orders, err := c.orderService.ListOrders(uint(userID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.NewAPIResponse(http.StatusNotFound, "No orders found", nil, "").Send(ctx)
			return
		}
		utils.NewAPIResponse(http.StatusInternalServerError, "Failed to retrieve orders", nil, err.Error()).Send(ctx)
		return
	}

	utils.NewAPIResponse(http.StatusOK, "Orders retrieved successfully", orders, "").Send(ctx)
}

// CancelOrder godoc
// @Summary      Cancel an order
// @Description  Allows a user to cancel a pending order
// @Tags         orders
// @Param        id   path      string  true  "Order ID"
// @Success      200  {object}  utils.APIResponse
// @Failure      400  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Security     BearerAuth
// @Router       /orders/{id}/cancel [put]
func (c *OrderController) CancelOrder(ctx *gin.Context) {
	orderID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.NewAPIResponse(http.StatusBadRequest, "Invalid order ID", nil, err.Error()).Send(ctx)
		return
	}

	userIDStr, exists := ctx.Get("userID")
	if !exists {
		utils.NewAPIResponse(http.StatusUnauthorized, "Unauthorized", nil, "User ID not found in context").Send(ctx)
		return
	}

	userID, err := strconv.ParseUint(userIDStr.(string), 10, 32)
	if err != nil {
		utils.NewAPIResponse(http.StatusBadRequest, "Invalid User ID", nil, "User ID conversion failed").Send(ctx)
		return
	}

	err = c.orderService.CancelOrder(uint(orderID), uint(userID))
	if err != nil {
		switch err.Error() {
		case "order not found":
			utils.NewAPIResponse(http.StatusNotFound, "Order not found", nil, "").Send(ctx)
		case "order is not eligible for cancellation":
			utils.NewAPIResponse(http.StatusBadRequest, "Order cannot be canceled", nil, "Order status must be 'Pending' to cancel").Send(ctx)
		default:
			utils.NewAPIResponse(http.StatusInternalServerError, "Failed to cancel order", nil, err.Error()).Send(ctx)
		}
		return
	}

	utils.NewAPIResponse(http.StatusOK, "Order canceled successfully", nil, "").Send(ctx)
}


// UpdateOrderStatus godoc
// @Summary      Update an order status
// @Description  Admin only: Updates the status of an order. The valid statuses are:
//               - pending: The order is newly created and awaiting processing.
//               - processing: The order is being processed.
//               - shipped: The order has been shipped to the customer.
//               - delivered: The order has been delivered to the customer.
//               - cancelled: The order was cancelled and will not be fulfilled.
// @Tags         orders
// @Param        id      path      string               true  "Order ID"
// @Param        status  body      UpdateOrderStatusDTO true  "New order status. Allowed values are 'pending', 'processing', 'shipped', 'delivered', 'cancelled'"
// @Success      200     {object}  utils.APIResponse{data=models.Order}
// @Failure      400     {object}  utils.APIResponse
// @Failure      404     {object}  utils.APIResponse
// @Failure      500     {object}  utils.APIResponse
// @Security     BearerAuth
// @Router       /orders/{id}/status [put]

func (c *OrderController) UpdateOrderStatus(ctx *gin.Context) {
	orderID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.NewAPIResponse(http.StatusBadRequest, "Invalid order ID", nil, err.Error()).Send(ctx)
		return
	}

	var input UpdateOrderStatusDTO
	if err := ctx.ShouldBindJSON(&input); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			for _, fieldErr := range ve {
				if fieldErr.Tag() == "orderStatus" {
					utils.NewAPIResponse(http.StatusBadRequest, "Invalid order status", nil, models.OrderStatus("").IsValid().Error()).Send(ctx)
					return
				}
			}
		}
		utils.NewAPIResponse(http.StatusBadRequest, "Invalid input", nil, err.Error()).Send(ctx)
		return
	}

	order, err := c.orderService.UpdateOrderStatus(uint(orderID), input.Status)
	if err != nil {
		if err.Error() == "order not found" {
			utils.NewAPIResponse(http.StatusNotFound, err.Error(), nil, "").Send(ctx)
		} else {
			utils.NewAPIResponse(http.StatusInternalServerError, "Failed to update order status", nil, err.Error()).Send(ctx)
		}
		return
	}

	utils.NewAPIResponse(http.StatusOK, "Order status updated successfully", order, "").Send(ctx)
}
