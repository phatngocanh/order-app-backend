package service

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/pna/order-app-backend/internal/domain/model"
)

type InventoryService interface {
	GetAll(ctx context.Context) (*model.GetAllInventoryResponse, string)
	GetByProductID(ctx *gin.Context, productID int) (*model.InventoryResponse, string)
	UpdateQuantity(ctx *gin.Context, productID int, request model.UpdateInventoryQuantityRequest) (*model.InventoryResponse, string)
}
