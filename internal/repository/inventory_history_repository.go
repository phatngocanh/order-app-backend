package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pna/order-app-backend/internal/domain/entity"
)

type InventoryHistoryRepository interface {
	GetAllByProductIDQuery(ctx context.Context, productID int, tx *sqlx.Tx) ([]entity.InventoryHistory, error)
	CreateCommand(ctx context.Context, inventoryHistory *entity.InventoryHistory, tx *sqlx.Tx) error
}
