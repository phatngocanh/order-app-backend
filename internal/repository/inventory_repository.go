package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pna/order-app-backend/internal/domain/entity"
)

type InventoryRepository interface {
	CreateCommand(ctx context.Context, inventory *entity.Inventory, tx *sqlx.Tx) error
	GetOneByProductIDQuery(ctx context.Context, productID int, tx *sqlx.Tx) (*entity.Inventory, error)
	UpdateQuantityCommand(ctx context.Context, productID int, quantity int, version string, tx *sqlx.Tx) error
}
