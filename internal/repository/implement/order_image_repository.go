package repositoryimplement

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/pna/order-app-backend/internal/database"
	"github.com/pna/order-app-backend/internal/domain/entity"
	"github.com/pna/order-app-backend/internal/repository"
	"github.com/pna/order-app-backend/internal/utils/error_utils"
)

type OrderImageRepository struct {
	db *sqlx.DB
}

func NewOrderImageRepository(db database.Db) repository.OrderImageRepository {
	return &OrderImageRepository{db: db}
}

func (repo *OrderImageRepository) GetAllByOrderIDQuery(ctx context.Context, orderID int, tx *sqlx.Tx) ([]entity.OrderImage, error) {
	var orderImages []entity.OrderImage
	query := "SELECT * FROM order_images WHERE order_id = ? ORDER BY id"
	var err error
	if tx != nil {
		err = tx.SelectContext(ctx, &orderImages, query, orderID)
	} else {
		err = repo.db.SelectContext(ctx, &orderImages, query, orderID)
	}
	if err != nil {
		return nil, err
	}

	// Ensure we always return an empty slice instead of nil
	if orderImages == nil {
		orderImages = make([]entity.OrderImage, 0)
	}

	return orderImages, nil
}

func (repo *OrderImageRepository) GetOneByIDQuery(ctx context.Context, id int, tx *sqlx.Tx) (*entity.OrderImage, error) {
	var orderImage entity.OrderImage
	query := "SELECT * FROM order_images WHERE id = ?"
	var err error
	if tx != nil {
		err = tx.GetContext(ctx, &orderImage, query, id)
	} else {
		err = repo.db.GetContext(ctx, &orderImage, query, id)
	}
	if err != nil {
		if err.Error() == error_utils.SystemErrorMessage.SqlxNoRow {
			return nil, nil
		}
		return nil, err
	}
	return &orderImage, nil
}

func (repo *OrderImageRepository) CreateCommand(ctx context.Context, orderImage *entity.OrderImage, tx *sqlx.Tx) error {
	insertQuery := `INSERT INTO order_images(order_id, image_url, s3_key) VALUES (:order_id, :image_url, :s3_key)`
	var result sql.Result
	var err error
	if tx != nil {
		result, err = tx.NamedExecContext(ctx, insertQuery, orderImage)
	} else {
		result, err = repo.db.NamedExecContext(ctx, insertQuery, orderImage)
	}
	if err != nil {
		return err
	}
	lastID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	orderImage.ID = int(lastID)
	return nil
}

func (repo *OrderImageRepository) DeleteByIDCommand(ctx context.Context, id int, tx *sqlx.Tx) error {
	deleteQuery := `DELETE FROM order_images WHERE id = ?`
	if tx != nil {
		_, err := tx.ExecContext(ctx, deleteQuery, id)
		return err
	}
	_, err := repo.db.ExecContext(ctx, deleteQuery, id)
	return err
}
