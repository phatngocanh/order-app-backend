package repositoryimplement

import (
	"context"
	"database/sql"
	"strings"

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

func (repo *InventoryRepository) GetAllQuery(ctx context.Context, tx *sqlx.Tx) ([]entity.Inventory, error) {
	var inventories []entity.Inventory
	query := "SELECT * FROM inventory ORDER BY id"
	var err error

	if tx != nil {
		err = tx.SelectContext(ctx, &inventories, query)
	} else {
		err = repo.db.SelectContext(ctx, &inventories, query)
	}

	if err != nil {
		return nil, err
	}

	return inventories, nil
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

	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, updateQuery, quantity, version, productID)
	} else {
		_, err = repo.db.ExecContext(ctx, updateQuery, quantity, version, productID)
	}

	if err != nil {
		// Check if it's a constraint violation error
		if strings.Contains(err.Error(), "check_quantity_non_negative") {
			return &error_utils.ConstraintViolationError{Message: "Số lượng kho không thể âm"}
		}
		return err
	}

	return nil
}

func (repo *InventoryRepository) GetOneByIDForUpdateQuery(ctx context.Context, productID int, tx *sqlx.Tx) (*entity.Inventory, error) {
	var inventory entity.Inventory
	query := "SELECT * FROM inventory WHERE id = ? FOR UPDATE"
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

func (repo *InventoryRepository) UpdateQuantityWithVersionCommand(ctx context.Context, productID int, quantity int, expectedVersion string, newVersion string, tx *sqlx.Tx) error {
	updateQuery := `UPDATE inventory SET quantity = quantity + ?, version = ? WHERE product_id = ? AND version = ?`

	var result sql.Result
	var err error
	if tx != nil {
		result, err = tx.ExecContext(ctx, updateQuery, quantity, newVersion, productID, expectedVersion)
	} else {
		result, err = repo.db.ExecContext(ctx, updateQuery, quantity, newVersion, productID, expectedVersion)
	}

	if err != nil {
		// Check if it's a constraint violation error
		if strings.Contains(err.Error(), "check_quantity_non_negative") {
			return &error_utils.ConstraintViolationError{Message: "Quantity cannot be negative"}
		}
		return err
	}

	// Check if any rows were affected (version mismatch)
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return &error_utils.VersionMismatchError{Message: "Version mismatch. Try again"}
	}

	return nil
}

func (repo *InventoryRepository) SelectManyForUpdate(ctx context.Context, ids []int, tx *sqlx.Tx) ([]entity.Inventory, error) {
	if len(ids) == 0 {
		return []entity.Inventory{}, nil
	}
	query, args, err := sqlx.In("SELECT * FROM inventory WHERE id IN (?) FOR UPDATE", ids)
	if err != nil {
		return nil, err
	}
	query = repo.db.Rebind(query)
	var inventories []entity.Inventory
	if tx != nil {
		err = tx.SelectContext(ctx, &inventories, query, args...)
	} else {
		err = repo.db.SelectContext(ctx, &inventories, query, args...)
	}
	if err != nil {
		return nil, err
	}
	return inventories, nil
}

func (repo *InventoryRepository) GetInventoryIDsByProductIDsQuery(ctx context.Context, productIDs []int, tx *sqlx.Tx) ([]int, error) {
	if len(productIDs) == 0 {
		return []int{}, nil
	}
	query, args, err := sqlx.In("SELECT id FROM inventory WHERE product_id IN (?)", productIDs)
	if err != nil {
		return nil, err
	}
	query = repo.db.Rebind(query)
	rows, err := repo.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}
