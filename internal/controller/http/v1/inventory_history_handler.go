package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	httpcommon "github.com/pna/order-app-backend/internal/domain/http_common"
	"github.com/pna/order-app-backend/internal/service"
	"github.com/pna/order-app-backend/internal/utils/error_utils"
)

type InventoryHistoryHandler struct {
	inventoryHistoryService service.InventoryHistoryService
}

func NewInventoryHistoryHandler(inventoryHistoryService service.InventoryHistoryService) *InventoryHistoryHandler {
	return &InventoryHistoryHandler{
		inventoryHistoryService: inventoryHistoryService,
	}
}

// @Summary Get All Inventory Histories
// @Description Retrieve all inventory history records for a specific product
// @Tags Inventory History
// @Produce json
// @Param  Authorization header string true "Authorization: Bearer"
// @Param productId path int true "Product ID"
// @Success 200 {object} httpcommon.HttpResponse[model.GetAllInventoryHistoriesResponse]
// @Failure 500 {object} httpcommon.HttpResponse[any]
// @Router /products/{productId}/inventories/histories [get]
func (h *InventoryHistoryHandler) GetAll(ctx *gin.Context) {
	productID, err := strconv.Atoi(ctx.Param("productId"))
	if err != nil {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "productId")
		ctx.JSON(statusCode, errResponse)
		return
	}
	response, errCode := h.inventoryHistoryService.GetAll(ctx, productID)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.JSON(http.StatusOK, httpcommon.NewSuccessResponse(response))
}
