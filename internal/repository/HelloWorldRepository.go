package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/pna/order-app-backend/internal/domain/entity"
)

type HelloWorldRepository interface {
	GetHelloWorldQuery(ctx context.Context, tx *sqlx.Tx) (entity.HelloWorld, error)
}
