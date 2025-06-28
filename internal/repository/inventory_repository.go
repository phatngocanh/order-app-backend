package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pna/order-app-backend/internal/domain/entity"
)

type InventoryRepository interface {
	CreateCommand(ctx context.Context, inventory *entity.Inventory, tx *sqlx.Tx) error
	GetAllQuery(ctx context.Context, tx *sqlx.Tx) ([]entity.Inventory, error)
	GetOneByProductIDQuery(ctx context.Context, productID int, tx *sqlx.Tx) (*entity.Inventory, error)
	UpdateQuantityCommand(ctx context.Context, productID int, quantity int, version string, tx *sqlx.Tx) error
	GetOneByIDForUpdateQuery(ctx context.Context, productID int, tx *sqlx.Tx) (*entity.Inventory, error)
	UpdateQuantityWithVersionCommand(ctx context.Context, productID int, quantity int, expectedVersion string, newVersion string, tx *sqlx.Tx) error

	SelectManyForUpdate(ctx context.Context, ids []int, tx *sqlx.Tx) ([]entity.Inventory, error)
	GetInventoryIDsByProductIDsQuery(ctx context.Context, productIDs []int, tx *sqlx.Tx) ([]int, error)
}
