package service

import (
	"context"

	"github.com/pna/order-app-backend/internal/domain/model"
)

type OrderService interface {
	GetAll(ctx context.Context) (model.GetAllOrdersResponse, error)
	GetOne(ctx context.Context, id int) (model.GetOneOrderResponse, error)
	Create(ctx context.Context, req model.CreateOrderRequest) error
	Update(ctx context.Context, req model.UpdateOrderRequest) error
}
