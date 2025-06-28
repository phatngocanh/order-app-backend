package service

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/pna/order-app-backend/internal/domain/model"
)

type OrderService interface {
	GetAll(ctx context.Context, userID int, customerID int, deliveryStatuses string, sortBy string) (model.GetAllOrdersResponse, string)
	GetOne(ctx context.Context, id int) (model.GetOneOrderResponse, string)
	Create(ctx *gin.Context, req model.CreateOrderRequest) string
	Update(ctx context.Context, req model.UpdateOrderRequest) string
	Delete(ctx *gin.Context, id int) string
}
