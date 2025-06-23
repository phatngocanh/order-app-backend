package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pna/order-app-backend/internal/domain/entity"
)

type CustomerRepository interface {
	GetAllQuery(ctx context.Context, tx *sqlx.Tx) ([]entity.Customer, error)
	GetOneByIDQuery(ctx context.Context, id int, tx *sqlx.Tx) (*entity.Customer, error)
	CreateCommand(ctx context.Context, customer *entity.Customer, tx *sqlx.Tx) error
	UpdateCommand(ctx context.Context, customer *entity.Customer, tx *sqlx.Tx) error
}
