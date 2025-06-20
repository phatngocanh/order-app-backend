package repositoryimplement

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/pna/order-app-backend/internal/database"
	"github.com/pna/order-app-backend/internal/domain/entity"
	"github.com/pna/order-app-backend/internal/repository"
)

type HelloWorldRepository struct {
	db *sqlx.DB
}

func NewHelloWorldRepository(db database.Db) repository.HelloWorldRepository {
	return &HelloWorldRepository{
		db: db,
	}
}

func (r HelloWorldRepository) GetHelloWorldQuery(ctx context.Context, tx *sqlx.Tx) (entity.HelloWorld, error) {
	return entity.HelloWorld{
		Message: "Hello World!",
	}, nil
}
