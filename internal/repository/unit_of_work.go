package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type UnitOfWork interface {
	Begin(ctx context.Context) (*sqlx.Tx, error)
	Commit(tx *sqlx.Tx) error
	Rollback(tx *sqlx.Tx) error
}
