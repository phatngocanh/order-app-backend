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

type CustomerHandler struct {
	customerService service.CustomerService
}

func NewCustomerHandler(customerService service.CustomerService) *CustomerHandler {
	return &CustomerHandler{
		customerService: customerService,
	}
}

// @Summary Create Customer
// @Description Create a new customer
// @Tags Customers
// @Accept json
// @Produce json
// @Param  Authorization header string true "Authorization: Bearer"
// @Param request body model.CreateCustomerRequest true "Customer information"
// @Success 201 {object} httpcommon.HttpResponse[model.CustomerResponse]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
// @Router /customers [post]
func (h *CustomerHandler) Create(ctx *gin.Context) {
	var request model.CreateCustomerRequest
	if err := validation.BindJsonAndValidate(ctx, &request); err != nil {
		return
	}

	response, errCode := h.customerService.Create(ctx, request)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.JSON(http.StatusCreated, httpcommon.NewSuccessResponse(response))
}

// @Summary Update Customer
// @Description Update an existing customer
// @Tags Customers
// @Accept json
// @Produce json
// @Param  Authorization header string true "Authorization: Bearer"
// @Param request body model.UpdateCustomerRequest true "Updated customer information"
// @Success 200 {object} httpcommon.HttpResponse[model.CustomerResponse]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 404 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
// @Router /customers/{customerId} [put]
func (h *CustomerHandler) Update(ctx *gin.Context) {
	customerID, err := strconv.Atoi(ctx.Param("customerId"))
	if err != nil {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "customerId")
		ctx.JSON(statusCode, errResponse)
		return
	}

	var request model.UpdateCustomerRequest
	if err := validation.BindJsonAndValidate(ctx, &request); err != nil {
		return
	}

	response, errCode := h.customerService.Update(ctx, customerID, request)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.JSON(http.StatusOK, httpcommon.NewSuccessResponse(response))
}

// @Summary Get All Customers
// @Description Retrieve all customers
// @Tags Customers
// @Produce json
// @Param  Authorization header string true "Authorization: Bearer"
// @Success 200 {object} httpcommon.HttpResponse[model.GetAllCustomersResponse]
// @Failure 500 {object} httpcommon.HttpResponse[any]
// @Router /customers [get]
func (h *CustomerHandler) GetAll(ctx *gin.Context) {
	response, errCode := h.customerService.GetAll(ctx)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.JSON(http.StatusOK, httpcommon.NewSuccessResponse(response))
}

// @Summary Get Customer by ID
// @Description Retrieve a customer by its ID
// @Tags Customers
// @Produce json
// @Param  Authorization header string true "Authorization: Bearer"
// @Param id path int true "Customer ID"
// @Success 200 {object} httpcommon.HttpResponse[model.GetOneCustomerResponse]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 404 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
// @Router /customers/{customerId} [get]
func (h *CustomerHandler) GetOne(ctx *gin.Context) {
	idStr := ctx.Param("customerId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "id")
		ctx.JSON(statusCode, errResponse)
		return
	}

	response, errCode := h.customerService.GetOne(ctx, id)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.JSON(http.StatusOK, httpcommon.NewSuccessResponse(response))
}
