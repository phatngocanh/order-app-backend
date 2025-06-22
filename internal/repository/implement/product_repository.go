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

type ProductRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db database.Db) repository.ProductRepository {
	return &ProductRepository{db: db}
}

func (repo *ProductRepository) GetAllQuery(ctx context.Context, tx *sqlx.Tx) ([]entity.Product, error) {
	var products []entity.Product
	query := "SELECT * FROM products ORDER BY id"
	var err error

	if tx != nil {
		err = tx.SelectContext(ctx, &products, query)
	} else {
		err = repo.db.SelectContext(ctx, &products, query)
	}

	if err != nil {
		return nil, err
	}

	if products == nil {
		return []entity.Product{}, nil
	}

	return products, nil
}

func (repo *ProductRepository) GetOneByIDQuery(ctx context.Context, id int, tx *sqlx.Tx) (*entity.Product, error) {
	var product entity.Product
	query := "SELECT * FROM products WHERE id = ?"
	var err error

	if tx != nil {
		err = tx.GetContext(ctx, &product, query, id)
	} else {
		err = repo.db.GetContext(ctx, &product, query, id)
	}

	if err != nil {
		if err.Error() == error_utils.SystemErrorMessage.SqlxNoRow {
			return nil, nil
		}
		return nil, err
	}

	return &product, nil
}

func (repo *ProductRepository) CreateCommand(ctx context.Context, product *entity.Product, tx *sqlx.Tx) error {
	insertQuery := `INSERT INTO products(name, spec, original_price) VALUES (:name, :spec, :original_price)`

	var result sql.Result
	var err error

	if tx != nil {
		result, err = tx.NamedExecContext(ctx, insertQuery, product)
	} else {
		result, err = repo.db.NamedExecContext(ctx, insertQuery, product)
	}

	if err != nil {
		return err
	}

	// Get the last inserted ID
	lastID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// Set the ID to the product entity
	product.ID = int(lastID)
	return nil
}

func (repo *ProductRepository) UpdateCommand(ctx context.Context, product *entity.Product, tx *sqlx.Tx) error {
	updateQuery := `UPDATE products SET name = :name, spec = :spec, original_price = :original_price WHERE id = :id`

	if tx != nil {
		_, err := tx.NamedExecContext(ctx, updateQuery, product)
		return err
	}
	_, err := repo.db.NamedExecContext(ctx, updateQuery, product)
	return err
}
