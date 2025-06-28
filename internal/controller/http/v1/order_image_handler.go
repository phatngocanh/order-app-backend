package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	httpcommon "github.com/pna/order-app-backend/internal/domain/http_common"
	"github.com/pna/order-app-backend/internal/service"
	"github.com/pna/order-app-backend/internal/utils/error_utils"
)

type OrderImageHandler struct {
	orderImageService service.OrderImageService
}

func NewOrderImageHandler(orderImageService service.OrderImageService) *OrderImageHandler {
	return &OrderImageHandler{
		orderImageService: orderImageService,
	}
}

// @Summary Generate Signed Upload URL
// @Description Generate a signed URL for uploading an image to S3
// @Tags Order Images
// @Produce json
// @Param  Authorization header string true "Authorization: Bearer"
// @Param orderId path int true "Order ID"
// @Param fileName query string true "File name"
// @Param contentType query string true "Content type (e.g., image/jpeg)"
// @Success 200 {object} httpcommon.HttpResponse[model.GenerateSignedUploadURLResponse]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 401 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
// @Router /orders/{orderId}/images/upload-url [post]
func (h *OrderImageHandler) GenerateSignedUploadURL(ctx *gin.Context) {
	// Get order ID from path parameter
	orderID, err := strconv.Atoi(ctx.Param("orderId"))
	if err != nil {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "orderId")
		ctx.JSON(statusCode, errResponse)
		return
	}

	// Get query parameters
	fileName := ctx.Query("fileName")
	if fileName == "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "fileName is required")
		ctx.JSON(statusCode, errResponse)
		return
	}

	contentType := ctx.Query("contentType")
	if contentType == "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "contentType is required")
		ctx.JSON(statusCode, errResponse)
		return
	}

	// Generate signed upload URL
	response, errCode := h.orderImageService.GenerateSignedUploadURL(ctx, orderID, fileName, contentType)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.JSON(http.StatusOK, httpcommon.NewSuccessResponse(&response))
}

// @Summary Upload Order Image
// @Description Upload an image for a specific order
// @Tags Order Images
// @Accept multipart/form-data
// @Produce json
// @Param  Authorization header string true "Authorization: Bearer"
// @Param orderId path int true "Order ID"
// @Param file formData file true "Image file to upload"
// @Success 201 {object} httpcommon.HttpResponse[model.UploadOrderImageResponse]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 401 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
// @Router /orders/{orderId}/images [post]
func (h *OrderImageHandler) UploadImage(ctx *gin.Context) {
	// Get order ID from path parameter
	orderID, err := strconv.Atoi(ctx.Param("orderId"))
	if err != nil {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "orderId")
		ctx.JSON(statusCode, errResponse)
		return
	}

	// Get the uploaded file
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "file")
		ctx.JSON(statusCode, errResponse)
		return
	}
	defer file.Close()

	// Validate file type (optional - you can add more validation)
	if header.Size == 0 {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "file is empty")
		ctx.JSON(statusCode, errResponse)
		return
	}

	// Upload the image
	response, errCode := h.orderImageService.UploadImage(ctx, orderID, file, header.Filename)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.JSON(http.StatusCreated, httpcommon.NewSuccessResponse(&response))
}

// @Summary Get Order Images
// @Description Get all images for a specific order
// @Tags Order Images
// @Produce json
// @Param  Authorization header string true "Authorization: Bearer"
// @Param orderId path int true "Order ID"
// @Success 200 {object} httpcommon.HttpResponse[model.GetOrderImagesResponse]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
// @Router /orders/{orderId}/images [get]
func (h *OrderImageHandler) GetImagesByOrderID(ctx *gin.Context) {
	// Get order ID from path parameter
	orderID, err := strconv.Atoi(ctx.Param("orderId"))
	if err != nil {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "orderId")
		ctx.JSON(statusCode, errResponse)
		return
	}

	// Get images for the order
	response, errCode := h.orderImageService.GetImagesByOrderID(ctx, orderID)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.JSON(http.StatusOK, httpcommon.NewSuccessResponse(&response))
}

// @Summary Delete Order Image
// @Description Delete a specific image from an order
// @Tags Order Images
// @Produce json
// @Param  Authorization header string true "Authorization: Bearer"
// @Param orderId path int true "Order ID"
// @Param imageId path int true "Image ID"
// @Success 200 {object} httpcommon.HttpResponse[any]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 401 {object} httpcommon.HttpResponse[any]
// @Failure 404 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
// @Router /orders/{orderId}/images/{imageId} [delete]
func (h *OrderImageHandler) DeleteImage(ctx *gin.Context) {
	// Get image ID from path parameter
	imageID, err := strconv.Atoi(ctx.Param("imageId"))
	if err != nil {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "imageId")
		ctx.JSON(statusCode, errResponse)
		return
	}

	// Delete the image
	errCode := h.orderImageService.DeleteImage(ctx, imageID)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.JSON(http.StatusOK, httpcommon.NewSuccessResponse[any](nil))
}
