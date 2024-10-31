package orders

import (
	"ecommerce-api/utils"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderController struct {
	orderService *OrderService
}

func NewOrderController(orderService *OrderService) *OrderController {
	return &OrderController{orderService: orderService}
}

// PlaceOrder godoc
// @Summary      Place an order
// @Description  Allows a user to place an order for one or more products
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        products  body      []PlaceOrderDTO   true  "List of products to order"
// @Success      201       {object}  utils.APIResponse{data=models.Order}
// @Failure      400       {object}  utils.APIResponse
// @Failure      404       {object}  utils.APIResponse
// @Failure      500       {object}  utils.APIResponse
// @Security     BearerAuth
// @Router       /api/v1/orders [post]
func (c *OrderController) PlaceOrder(ctx *gin.Context) {
	var input PlaceOrderDTO

	// Bind JSON input to struct
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.NewAPIResponse(http.StatusBadRequest, "Invalid input", nil, err.Error()).Send(ctx)
		return
	}

	// Get user ID from context
	userID := ctx.GetUint("userID")

	// Place the order using the service
	order, err := c.orderService.PlaceOrder(userID, input.Products)
	if err != nil {
		utils.NewAPIResponse(http.StatusInternalServerError, "Failed to place order", nil, err.Error()).Send(ctx)
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
// @Router       /api/v1/orders [get]
func (c *OrderController) ListOrders(ctx *gin.Context) {
	// Get user ID from context
	userID := ctx.GetUint("userID")

	orders, err := c.orderService.ListOrders(userID)
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
// @Router       /api/v1/orders/{id}/cancel [put]
func (c *OrderController) CancelOrder(ctx *gin.Context) {
	orderID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.NewAPIResponse(http.StatusBadRequest, "Invalid order ID", nil, err.Error()).Send(ctx)
		return
	}

	userID := ctx.GetUint("userID")
	err = c.orderService.CancelOrder(uint(orderID), userID)
	if err != nil {
		if err.Error() == "order not found or not eligible for cancellation" {
			utils.NewAPIResponse(http.StatusNotFound, err.Error(), nil, "").Send(ctx)
		} else {
			utils.NewAPIResponse(http.StatusInternalServerError, "Failed to cancel order", nil, err.Error()).Send(ctx)
		}
		return
	}

	utils.NewAPIResponse(http.StatusOK, "Order canceled successfully", nil, "").Send(ctx)
}

// UpdateOrderStatus godoc
// @Summary      Update an order status
// @Description  Admin only: Updates the status of an order
// @Tags         orders
// @Param        id      path      string               true  "Order ID"
// @Param        status  body      UpdateOrderStatusDTO   true  "New order status"
// @Success      200     {object}  utils.APIResponse{data=models.Order}
// @Failure      400     {object}  utils.APIResponse
// @Failure      404     {object}  utils.APIResponse
// @Failure      500     {object}  utils.APIResponse
// @Security     BearerAuth
// @Router       /api/v1/orders/{id}/status [put]
func (c *OrderController) UpdateOrderStatus(ctx *gin.Context) {
	orderID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.NewAPIResponse(http.StatusBadRequest, "Invalid order ID", nil, err.Error()).Send(ctx)
		return
	}

	var input UpdateOrderStatusDTO
	if err := ctx.ShouldBindJSON(&input); err != nil {
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
