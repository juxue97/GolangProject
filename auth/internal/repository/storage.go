package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/juxue97/auth/internal/db"
	"github.com/juxue97/common"
)

var QueryTimeoutDuration = time.Second * 5

type Repository struct {
	Users interface {
		Create(context.Context, *User) error
		GetByID(context.Context, int64) (*User, error)
		GetByEmail(context.Context, string) (*User, error)
		CreateAndInvite(context.Context, *User, string, time.Duration) error
		ActivateUser(context.Context, string) error
		Delete(context.Context, int64) error
	}
}

func NewRepository(db *sql.DB) *Repository {
	if db == nil {
		common.Logger.Fatal("PgDB is nil")
	}
	return &Repository{
		Users: &UserStore{DB: db},
	}
}

var Store *Repository

func init() {
	Store = NewRepository(db.PgDB)
	common.Logger.Info("Store initialized")
}

func withTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
