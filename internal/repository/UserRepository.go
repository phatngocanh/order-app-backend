package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pna/order-app-backend/internal/domain/entity"
)

type UserRepository interface {
	CreateCommand(ctx context.Context, user *entity.User, tx *sqlx.Tx) error
}
