package users

import (
	"context"
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64    `json:"id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Password  password `json:"password"`
	CreatedAt string   `json:"created_at"`
	IsActive  bool     `json:"is_active"`
}

type password struct {
	text *string
	hash []byte
}

func (p *password) SetPassword(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.text = &text
	p.hash = hash

	return nil
}

func (p *password) ComparePassword(text string) bool {
	return bcrypt.CompareHashAndPassword(p.hash, []byte(text)) == nil
}

type UserStore struct {
	DB *sql.DB
}

func (us *UserStore) Create(ctx context.Context, user *User) error {
	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3)`

	err := us.DB.QueryRowContext(
		ctx,
		query, user.Username,
		user.Email,
		user.Password,
	).Scan(
		&user.ID,
		&user.CreatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (us *UserStore) GetByID(ctx context.Context, id int64) (*User, error) {
	return nil, nil
}

func (us *UserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	return nil, nil
}
