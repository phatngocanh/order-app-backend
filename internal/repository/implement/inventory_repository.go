package repositoryimplement

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pna/order-app-backend/internal/database"
	"github.com/pna/order-app-backend/internal/domain/entity"
	"github.com/pna/order-app-backend/internal/repository"
	"github.com/pna/order-app-backend/internal/utils/error_utils"
)

type InventoryRepository struct {
	db *sqlx.DB
}

func NewInventoryRepository(db database.Db) repository.InventoryRepository {
	return &InventoryRepository{db: db}
}

func (repo *InventoryRepository) CreateCommand(ctx context.Context, inventory *entity.Inventory, tx *sqlx.Tx) error {
	insertQuery := `INSERT INTO inventory(product_id, quantity, version) VALUES (:product_id, :quantity, :version)`

	if tx != nil {
		_, err := tx.NamedExecContext(ctx, insertQuery, inventory)
		return err
	}
	_, err := repo.db.NamedExecContext(ctx, insertQuery, inventory)
	return err
}

func (repo *InventoryRepository) GetOneByProductIDQuery(ctx context.Context, productID int, tx *sqlx.Tx) (*entity.Inventory, error) {
	var inventory entity.Inventory
	query := "SELECT * FROM inventory WHERE product_id = ?"
	var err error

	if tx != nil {
		err = tx.GetContext(ctx, &inventory, query, productID)
	} else {
		err = repo.db.GetContext(ctx, &inventory, query, productID)
	}

	if err != nil {
		if err.Error() == error_utils.SystemErrorMessage.SqlxNoRow {
			return nil, nil
		}
		return nil, err
	}

	return &inventory, nil
}

func (repo *InventoryRepository) UpdateQuantityCommand(ctx context.Context, productID int, quantity int, version string, tx *sqlx.Tx) error {
	updateQuery := `UPDATE inventory SET quantity = quantity + ?, version = ? WHERE product_id = ?`

	if tx != nil {
		_, err := tx.ExecContext(ctx, updateQuery, quantity, version, productID)
		return err
	}
	_, err := repo.db.ExecContext(ctx, updateQuery, quantity, version, productID)
	return err
}
