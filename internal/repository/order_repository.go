package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pna/order-app-backend/internal/domain/entity"
)

type OrderRepository interface {
	GetAllQuery(ctx context.Context, tx *sqlx.Tx) ([]entity.Order, error)
	GetAllWithFiltersQuery(ctx context.Context, customerID int, deliveryStatuses string, sortBy string, tx *sqlx.Tx) ([]entity.Order, error)
	GetOneByIDQuery(ctx context.Context, id int, tx *sqlx.Tx) (*entity.Order, error)
	CreateCommand(ctx context.Context, order *entity.Order, tx *sqlx.Tx) error
	UpdateCommand(ctx context.Context, order *entity.Order, tx *sqlx.Tx) error
	DeleteByIDCommand(ctx context.Context, id int, tx *sqlx.Tx) error
}
