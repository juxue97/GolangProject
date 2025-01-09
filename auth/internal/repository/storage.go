package repository

import (
	"context"
	"database/sql"

	users "github.com/juxue97/auth/internal/repository/users"
)

type Repository struct {
	Users interface {
		Create(context.Context, *users.User) error
		GetByID(context.Context, int64) (*users.User, error)
		GetByEmail(context.Context, string) (*users.User, error)
	}
}

func NewRepository(db *sql.DB) Repository {
	return Repository{
		Users: &users.UserStore{DB: db},
	}
}
