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

type OrderRepository struct {
	db *sqlx.DB
}

func NewOrderRepository(db database.Db) repository.OrderRepository {
	return &OrderRepository{db: db}
}

func (repo *OrderRepository) GetAllQuery(ctx context.Context, tx *sqlx.Tx) ([]entity.Order, error) {
	var orders []entity.Order
	query := "SELECT * FROM orders ORDER BY id"
	var err error
	if tx != nil {
		err = tx.SelectContext(ctx, &orders, query)
	} else {
		err = repo.db.SelectContext(ctx, &orders, query)
	}
	if err != nil {
		return nil, err
	}
	if orders == nil {
		return []entity.Order{}, nil
	}
	return orders, nil
}

func (repo *OrderRepository) GetAllWithFiltersQuery(ctx context.Context, customerID int, deliveryStatuses string, sortBy string, tx *sqlx.Tx) ([]entity.Order, error) {
	var orders []entity.Order
	query := "SELECT * FROM orders WHERE 1=1"
	var args []interface{}

	// Add customer filter
	if customerID > 0 {
		query += " AND customer_id = ?"
		args = append(args, customerID)
	}

	// Add delivery statuses filter using IN query
	if deliveryStatuses != "" {
		// Split the comma-separated statuses
		statusList := strings.Split(deliveryStatuses, ",")
		if len(statusList) > 0 {
			// Create placeholders for IN query
			placeholders := make([]string, len(statusList))
			for i := range statusList {
				placeholders[i] = "?"
				args = append(args, strings.TrimSpace(statusList[i]))
			}
			query += " AND delivery_status IN (" + strings.Join(placeholders, ",") + ")"
		}
	}

	// Add sorting
	switch sortBy {
	case "order_date_asc":
		query += " ORDER BY order_date ASC"
	case "order_date_desc":
		query += " ORDER BY order_date DESC"
	default:
		query += " ORDER BY id DESC"
	}

	var err error
	if tx != nil {
		err = tx.SelectContext(ctx, &orders, query, args...)
	} else {
		err = repo.db.SelectContext(ctx, &orders, query, args...)
	}
	if err != nil {
		return nil, err
	}
	if orders == nil {
		return []entity.Order{}, nil
	}
	return orders, nil
}

func (repo *OrderRepository) GetOneByIDQuery(ctx context.Context, id int, tx *sqlx.Tx) (*entity.Order, error) {
	var order entity.Order
	query := "SELECT * FROM orders WHERE id = ?"
	var err error
	if tx != nil {
		err = tx.GetContext(ctx, &order, query, id)
	} else {
		err = repo.db.GetContext(ctx, &order, query, id)
	}
	if err != nil {
		if err.Error() == error_utils.SystemErrorMessage.SqlxNoRow {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}

func (repo *OrderRepository) CreateCommand(ctx context.Context, order *entity.Order, tx *sqlx.Tx) error {
	insertQuery := `INSERT INTO orders(customer_id, order_date, delivery_status, debt_status, status_transitioned_at, total_original_cost, total_sales_revenue, additional_cost, additonal_cost_note) VALUES (:customer_id, :order_date, :delivery_status, :debt_status, :status_transitioned_at, :total_original_cost, :total_sales_revenue, :additional_cost, :additonal_cost_note)`
	var result sql.Result
	var err error
	if tx != nil {
		result, err = tx.NamedExecContext(ctx, insertQuery, order)
	} else {
		result, err = repo.db.NamedExecContext(ctx, insertQuery, order)
	}
	if err != nil {
		return err
	}
	lastID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	order.ID = int(lastID)
	return nil
}

func (repo *OrderRepository) UpdateCommand(ctx context.Context, order *entity.Order, tx *sqlx.Tx) error {
	updateQuery := `UPDATE orders SET customer_id = :customer_id, order_date = :order_date, delivery_status = :delivery_status, debt_status = :debt_status, status_transitioned_at = :status_transitioned_at, total_original_cost = :total_original_cost, total_sales_revenue = :total_sales_revenue, additional_cost = :additional_cost, additonal_cost_note = :additonal_cost_note WHERE id = :id`
	if tx != nil {
		_, err := tx.NamedExecContext(ctx, updateQuery, order)
		return err
	}
	_, err := repo.db.NamedExecContext(ctx, updateQuery, order)
	return err
}

func (repo *OrderRepository) DeleteByIDCommand(ctx context.Context, id int, tx *sqlx.Tx) error {
	deleteQuery := `DELETE FROM orders WHERE id = ?`
	if tx != nil {
		_, err := tx.ExecContext(ctx, deleteQuery, id)
		return err
	}
	_, err := repo.db.ExecContext(ctx, deleteQuery, id)
	return err
}
