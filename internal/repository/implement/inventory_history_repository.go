package repositoryimplement

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/pna/order-app-backend/internal/database"
	"github.com/pna/order-app-backend/internal/domain/entity"
	"github.com/pna/order-app-backend/internal/repository"
)

type InventoryHistoryRepository struct {
	db *sqlx.DB
}

func NewInventoryHistoryRepository(db database.Db) repository.InventoryHistoryRepository {
	return &InventoryHistoryRepository{db: db}
}

func (repo *InventoryHistoryRepository) GetAllByProductIDQuery(ctx context.Context, productID int, tx *sqlx.Tx) ([]entity.InventoryHistory, error) {
	var inventoryHistories []entity.InventoryHistory
	query := "SELECT * FROM inventory_histories WHERE product_id = ? ORDER BY imported_at DESC"
	var err error

	if tx != nil {
		err = tx.SelectContext(ctx, &inventoryHistories, query, productID)
	} else {
		err = repo.db.SelectContext(ctx, &inventoryHistories, query, productID)
	}

	if err != nil {
		return nil, err
	}

	if inventoryHistories == nil {
		return []entity.InventoryHistory{}, nil
	}

	return inventoryHistories, nil
}

func (repo *InventoryHistoryRepository) CreateCommand(ctx context.Context, inventoryHistory *entity.InventoryHistory, tx *sqlx.Tx) error {
	insertQuery := `INSERT INTO inventory_histories(product_id, quantity, final_quantity, importer_name, imported_at, note, reference_id) VALUES (:product_id, :quantity, :final_quantity, :importer_name, :imported_at, :note, :reference_id)`

	var result sql.Result
	var err error

	if tx != nil {
		result, err = tx.NamedExecContext(ctx, insertQuery, inventoryHistory)
	} else {
		result, err = repo.db.NamedExecContext(ctx, insertQuery, inventoryHistory)
	}

	if err != nil {
		return err
	}

	// Get the last inserted ID
	lastID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// Set the ID to the inventory history entity
	inventoryHistory.ID = int(lastID)
	return nil
}
