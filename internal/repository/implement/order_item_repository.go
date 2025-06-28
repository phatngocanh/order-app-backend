package repositoryimplement

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/pna/order-app-backend/internal/database"
	"github.com/pna/order-app-backend/internal/domain/entity"
	"github.com/pna/order-app-backend/internal/repository"
)

type OrderItemRepository struct {
	db *sqlx.DB
}

func NewOrderItemRepository(db database.Db) repository.OrderItemRepository {
	return &OrderItemRepository{db: db}
}

func (repo *OrderItemRepository) GetAllByOrderIDQuery(ctx context.Context, orderID int, tx *sqlx.Tx) ([]entity.OrderItem, error) {
	var orderItems []entity.OrderItem
	query := "SELECT * FROM order_items WHERE order_id = ? ORDER BY id"
	var err error
	if tx != nil {
		err = tx.SelectContext(ctx, &orderItems, query, orderID)
	} else {
		err = repo.db.SelectContext(ctx, &orderItems, query, orderID)
	}
	if err != nil {
		return nil, err
	}
	if orderItems == nil {
		return []entity.OrderItem{}, nil
	}
	return orderItems, nil
}

func (repo *OrderItemRepository) CreateCommand(ctx context.Context, orderItem *entity.OrderItem, tx *sqlx.Tx) error {
	insertQuery := `INSERT INTO order_items(order_id, product_id, number_of_boxes, spec, quantity, selling_price, original_price, discount, final_amount, export_from) VALUES (:order_id, :product_id, :number_of_boxes, :spec, :quantity, :selling_price, :original_price, :discount, :final_amount, :export_from)`
	var result sql.Result
	var err error
	if tx != nil {
		result, err = tx.NamedExecContext(ctx, insertQuery, orderItem)
	} else {
		result, err = repo.db.NamedExecContext(ctx, insertQuery, orderItem)
	}
	if err != nil {
		return err
	}
	lastID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	orderItem.ID = int(lastID)
	return nil
}

func (repo *OrderItemRepository) DeleteByIDCommand(ctx context.Context, id int, tx *sqlx.Tx) error {
	deleteQuery := `DELETE FROM order_items WHERE id = ?`
	if tx != nil {
		_, err := tx.ExecContext(ctx, deleteQuery, id)
		return err
	}
	_, err := repo.db.ExecContext(ctx, deleteQuery, id)
	return err
}
