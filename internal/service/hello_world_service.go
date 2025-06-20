package service

import (
	"context"
	"github.com/pna/order-app-backend/internal/domain/model"
)

type HelloWorldService interface {
	HelloWorld(ctx context.Context) (*model.HelloWorldResponse, string)
}
