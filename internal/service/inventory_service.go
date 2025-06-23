package service

import (
	"github.com/gin-gonic/gin"
	"github.com/pna/order-app-backend/internal/domain/model"
)

type InventoryService interface {
	GetByProductID(ctx *gin.Context, productID int) (*model.InventoryResponse, string)
	UpdateQuantity(ctx *gin.Context, productID int, request model.UpdateInventoryQuantityRequest) (*model.InventoryResponse, string)
}
