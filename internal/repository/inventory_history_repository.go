package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pna/order-app-backend/internal/domain/entity"
)

type InventoryHistoryRepository interface {
	GetAllQuery(ctx context.Context, tx *sqlx.Tx) ([]entity.InventoryHistory, error)
	CreateCommand(ctx context.Context, inventoryHistory *entity.InventoryHistory, tx *sqlx.Tx) error
}
