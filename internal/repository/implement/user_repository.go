package repositoryimplement

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pna/order-app-backend/internal/database"
	"github.com/pna/order-app-backend/internal/domain/entity"
	"github.com/pna/order-app-backend/internal/repository"
	"github.com/pna/order-app-backend/internal/utils/error_utils"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db database.Db) repository.UserRepository {
	return &UserRepository{db: db}
}

func (repo *UserRepository) CreateCommand(ctx context.Context, user *entity.User, tx *sqlx.Tx) error {
	insertQuery := `INSERT INTO users(username, password) VALUES (:username, :password)`
	if tx != nil {
		_, err := tx.NamedExecContext(ctx, insertQuery, user)
		return err
	}
	_, err := repo.db.NamedExecContext(ctx, insertQuery, user)
	return err
}

func (repo *UserRepository) FindByUsernameQuery(ctx context.Context, username string, tx *sqlx.Tx) (*entity.User, error) {
	var user entity.User
	query := "SELECT * FROM users WHERE username = ?"
	var err error
	if tx != nil {
		err = tx.GetContext(ctx, &user, query, username)
	} else {
		err = repo.db.GetContext(ctx, &user, query, username)
	}

	if err != nil {
		if err.Error() == error_utils.SystemErrorMessage.SqlxNoRow {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (repo *UserRepository) FindByIDQuery(ctx context.Context, id int, tx *sqlx.Tx) (*entity.User, error) {
	var user entity.User
	query := "SELECT * FROM users WHERE id = ?"
	var err error
	if tx != nil {
		err = tx.GetContext(ctx, &user, query, id)
	} else {
		err = repo.db.GetContext(ctx, &user, query, id)
	}

	if err != nil {
		if err.Error() == error_utils.SystemErrorMessage.SqlxNoRow {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
