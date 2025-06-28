package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	httpcommon "github.com/pna/order-app-backend/internal/domain/http_common"
	"github.com/pna/order-app-backend/internal/domain/model"
	"github.com/pna/order-app-backend/internal/service"
	"github.com/pna/order-app-backend/internal/utils/error_utils"
	"github.com/pna/order-app-backend/internal/utils/validation"
)

type OrderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// @Summary Create Order
// @Description Create a new order with order items
// @Tags Orders
// @Accept json
// @Produce json
// @Param  Authorization header string true "Authorization: Bearer"
// @Param request body model.CreateOrderRequest true "Order information with items"
// @Success 201 {object} httpcommon.HttpResponse[any]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 401 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
// @Router /orders [post]
func (h *OrderHandler) Create(ctx *gin.Context) {
	var request model.CreateOrderRequest
	if err := validation.BindJsonAndValidate(ctx, &request); err != nil {
		return
	}

	errCode := h.orderService.Create(ctx, request)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.JSON(http.StatusCreated, httpcommon.NewSuccessResponse[any](nil))
}

// @Summary Update Order
// @Description Update an existing order
// @Tags Orders
// @Accept json
// @Produce json
// @Param  Authorization header string true "Authorization: Bearer"
// @Param request body model.UpdateOrderRequest true "Updated order information"
// @Success 200 {object} httpcommon.HttpResponse[any]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 404 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
// @Router /orders/{orderId} [put]
func (h *OrderHandler) Update(ctx *gin.Context) {
	orderID, err := strconv.Atoi(ctx.Param("orderId"))
	if err != nil {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "orderId")
		ctx.JSON(statusCode, errResponse)
		return
	}

	var request model.UpdateOrderRequest
	if err := validation.BindJsonAndValidate(ctx, &request); err != nil {
		return
	}

	request.ID = orderID
	errCode := h.orderService.Update(ctx, request)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.JSON(http.StatusOK, httpcommon.NewSuccessResponse[any](nil))
}

// @Summary Get All Orders
// @Description Retrieve all orders with optional filters and sorting
// @Tags Orders
// @Produce json
// @Param  Authorization header string true "Authorization: Bearer"
// @Param customer_id query int false "Filter by customer ID"
// @Param delivery_statuses query string false "Filter by delivery statuses (comma-separated, e.g., PENDING,DELIVERED)"
// @Param sort_by query string false "Sort by: order_date_asc, order_date_desc (default: id DESC)"
// @Success 200 {object} httpcommon.HttpResponse[model.GetAllOrdersResponse]
// @Failure 500 {object} httpcommon.HttpResponse[any]
// @Router /orders [get]
func (h *OrderHandler) GetAll(ctx *gin.Context) {
	// Get query parameters
	customerIDStr := ctx.Query("customer_id")
	deliveryStatuses := ctx.Query("delivery_statuses")
	sortBy := ctx.Query("sort_by")

	// Parse customer ID if provided
	customerID := 0
	if customerIDStr != "" {
		if id, err := strconv.Atoi(customerIDStr); err == nil {
			customerID = id
		}
	}

	response, errCode := h.orderService.GetAll(ctx, 0, customerID, deliveryStatuses, sortBy)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.JSON(http.StatusOK, httpcommon.NewSuccessResponse(&response))
}

// @Summary Get Order by ID
// @Description Retrieve an order by its ID
// @Tags Orders
// @Produce json
// @Param  Authorization header string true "Authorization: Bearer"
// @Param id path int true "Order ID"
// @Success 200 {object} httpcommon.HttpResponse[model.GetOneOrderResponse]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 404 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
// @Router /orders/{orderId} [get]
func (h *OrderHandler) GetOne(ctx *gin.Context) {
	idStr := ctx.Param("orderId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "id")
		ctx.JSON(statusCode, errResponse)
		return
	}

	response, errCode := h.orderService.GetOne(ctx, id)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.JSON(http.StatusOK, httpcommon.NewSuccessResponse(&response))
}

// @Summary Delete Order
// @Description Delete an order by its ID. If the order contains items exported from inventory, they will be restored to inventory.
// @Tags Orders
// @Produce json
// @Param  Authorization header string true "Authorization: Bearer"
// @Param orderId path int true "Order ID"
// @Success 200 {object} httpcommon.HttpResponse[any]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 401 {object} httpcommon.HttpResponse[any]
// @Failure 404 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
// @Router /orders/{orderId} [delete]
func (h *OrderHandler) Delete(ctx *gin.Context) {
	orderID, err := strconv.Atoi(ctx.Param("orderId"))
	if err != nil {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "orderId")
		ctx.JSON(statusCode, errResponse)
		return
	}

	errCode := h.orderService.Delete(ctx, orderID)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.JSON(http.StatusOK, httpcommon.NewSuccessResponse[any](nil))
}
