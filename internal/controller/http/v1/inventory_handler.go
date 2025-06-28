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

type InventoryHandler struct {
	inventoryService service.InventoryService
}

func NewInventoryHandler(inventoryService service.InventoryService) *InventoryHandler {
	return &InventoryHandler{
		inventoryService: inventoryService,
	}
}

// @Summary Get All Inventory
// @Description Retrieve all inventory information with product details
// @Tags Inventory
// @Produce json
// @Param  Authorization header string true "Authorization: Bearer"
// @Success 200 {object} httpcommon.HttpResponse[model.GetAllInventoryResponse]
// @Failure 500 {object} httpcommon.HttpResponse[any]
// @Router /inventory [get]
func (h *InventoryHandler) GetAll(ctx *gin.Context) {
	context := ctx.Request.Context()

	response, errCode := h.inventoryService.GetAll(context)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.JSON(http.StatusOK, httpcommon.NewSuccessResponse(response))
}

// @Summary Get Inventory by Product ID
// @Description Retrieve inventory information for a specific product
// @Tags Inventory
// @Produce json
// @Param  Authorization header string true "Authorization: Bearer"
// @Param productId path int true "Product ID"
// @Success 200 {object} httpcommon.HttpResponse[model.InventoryResponse]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 404 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
// @Router /products/{productId}/inventories [get]
func (h *InventoryHandler) GetByProductID(ctx *gin.Context) {
	productIDStr := ctx.Param("productId")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "productId")
		ctx.JSON(statusCode, errResponse)
		return
	}

	response, errCode := h.inventoryService.GetByProductID(ctx, productID)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.JSON(http.StatusOK, httpcommon.NewSuccessResponse(response))
}

// @Summary Update Inventory Quantity
// @Description Update the quantity of a product in inventory
// @Tags Inventory
// @Accept json
// @Produce json
// @Param  Authorization header string true "Authorization: Bearer"
// @Param productId path int true "Product ID"
// @Param request body model.UpdateInventoryQuantityRequest true "Quantity update information"
// @Success 200 {object} httpcommon.HttpResponse[model.InventoryResponse]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 404 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
// @Router /products/{productId}/inventories/quantity [put]
func (h *InventoryHandler) UpdateQuantity(ctx *gin.Context) {
	productIDStr := ctx.Param("productId")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "productId")
		ctx.JSON(statusCode, errResponse)
		return
	}

	var request model.UpdateInventoryQuantityRequest
	if err := validation.BindJsonAndValidate(ctx, &request); err != nil {
		return
	}

	response, errCode := h.inventoryService.UpdateQuantity(ctx, productID, request)
	if errCode != "" {
		field := ""
		if errCode == error_utils.ErrorCode.BAD_REQUEST {
			// Check if it's a version mismatch or quantity constraint error
			// For now, we'll use a generic approach since the service doesn't return specific field info
			field = "quantity" // Default to quantity field
		}
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, field)
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.JSON(http.StatusOK, httpcommon.NewSuccessResponse(response))
}
