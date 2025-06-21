package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pna/order-app-backend/internal/domain/entity"
)

type ProductRepository interface {
	GetAllQuery(ctx context.Context, tx *sqlx.Tx) ([]entity.Product, error)
	GetOneByIDQuery(ctx context.Context, id int, tx *sqlx.Tx) (*entity.Product, error)
	CreateCommand(ctx context.Context, product *entity.Product, tx *sqlx.Tx) error
	UpdateCommand(ctx context.Context, product *entity.Product, tx *sqlx.Tx) error
}
