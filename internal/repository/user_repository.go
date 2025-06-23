package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pna/order-app-backend/internal/domain/entity"
)

type UserRepository interface {
	CreateCommand(ctx context.Context, user *entity.User, tx *sqlx.Tx) error
	FindByUsernameQuery(ctx context.Context, username string, tx *sqlx.Tx) (*entity.User, error)
	FindByIDQuery(ctx context.Context, id int, tx *sqlx.Tx) (*entity.User, error)
}
