package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pna/order-app-backend/internal/domain/entity"
)

type OrderImageRepository interface {
	GetAllByOrderIDQuery(ctx context.Context, orderID int, tx *sqlx.Tx) ([]entity.OrderImage, error)
	GetOneByIDQuery(ctx context.Context, id int, tx *sqlx.Tx) (*entity.OrderImage, error)
	CreateCommand(ctx context.Context, orderImage *entity.OrderImage, tx *sqlx.Tx) error
	DeleteByIDCommand(ctx context.Context, id int, tx *sqlx.Tx) error
}
