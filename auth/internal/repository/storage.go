package repository

import (
	"context"
	"database/sql"
	"log"
	"time"
)

var QueryTimeoutDuration = time.Second * 5

type Repository struct {
	Users interface {
		Create(context.Context, *User) error
		GetAll(context.Context) ([]User, error)
		GetByID(context.Context, int64) (*User, error)
		GetByEmail(context.Context, string) (*User, error)
		Update(context.Context, *User) error
		DeleteUser(context.Context, int64) error
		CreateAndInvite(context.Context, *User, string, time.Duration) error
		ActivateUser(context.Context, string) error
		Delete(context.Context, int64) error
	}
	Roles interface {
		GetByName(context.Context, string) (*Role, error)
	}
}

func NewRepository(db *sql.DB) *Repository {
	if db == nil {
		log.Fatal("PgDB is nil")
	}
	return &Repository{
		Users: &UserStore{DB: db},
		Roles: &RoleStore{DB: db},
	}
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
