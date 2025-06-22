package repositoryimplement

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/pna/order-app-backend/internal/database"
	"github.com/pna/order-app-backend/internal/repository"
)

type UnitOfWorkImpl struct {
	db *sqlx.DB
}

func NewUnitOfWork(db database.Db) repository.UnitOfWork {
	return &UnitOfWorkImpl{db: db}
}

func (uow *UnitOfWorkImpl) Begin(ctx context.Context) (*sqlx.Tx, error) {
	return uow.db.BeginTxx(ctx, &sql.TxOptions{})
}

func (uow *UnitOfWorkImpl) Commit(tx *sqlx.Tx) error {
	if tx != nil {
		return tx.Commit()
	}
	return nil
}

func (uow *UnitOfWorkImpl) Rollback(tx *sqlx.Tx) error {
	if tx != nil {
		return tx.Rollback()
	}
	return nil
}
