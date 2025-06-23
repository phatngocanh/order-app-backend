package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pna/order-app-backend/internal/domain/entity"
)

type OrderItemRepository interface {
	GetAllByOrderIDQuery(ctx context.Context, orderID int, tx *sqlx.Tx) ([]entity.OrderItem, error)
	CreateCommand(ctx context.Context, orderItem *entity.OrderItem, tx *sqlx.Tx) error
	DeleteByIDCommand(ctx context.Context, id int, tx *sqlx.Tx) error
}
