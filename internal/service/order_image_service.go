package service

import (
	"github.com/gin-gonic/gin"
	"github.com/pna/order-app-backend/internal/domain/model"
)

type OrderImageService interface {
	DeleteImage(ctx *gin.Context, imageID int) string
	GenerateSignedUploadURL(ctx *gin.Context, orderID int, fileName string, contentType string) (model.GenerateSignedUploadURLResponse, string)
}
