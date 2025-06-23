package service

import (
	"context"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/pna/order-app-backend/internal/domain/model"
)

type OrderImageService interface {
	UploadImage(ctx *gin.Context, orderID int, file io.Reader, fileName string) (model.UploadOrderImageResponse, string)
	GetImagesByOrderID(ctx context.Context, orderID int) (model.GetOrderImagesResponse, string)
	DeleteImage(ctx *gin.Context, imageID int) string
}
