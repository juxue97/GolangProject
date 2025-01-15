package repository

import (
	"context"
	"time"

	"github.com/juxue97/common"
)

type (
	mockUserStore struct{}
	mockRoleStore struct{}
)

func NewMockStore() *Repository {
	return &Repository{
		Users: &mockUserStore{},
		Roles: &mockRoleStore{},
	}
}

func (rs *mockRoleStore) GetByName(context.Context, string) (*Role, error) {
	return nil, nil
}

func (ms *mockUserStore) Create(context.Context, *User) error {
	return nil
}

func (ms *mockUserStore) GetAll(context.Context) ([]User, error) {
	return nil, nil
}

func (ms *mockUserStore) GetByID(context.Context, int64) (*User, error) {
	return nil, nil
}

func (ms *mockUserStore) GetByEmail(context.Context, string) (*User, error) {
	return &User{
		ID:        22,
		Username:  "MehNohNah",
		Email:     "hwteh1997@gmail.com",
		Password:  password{hash: []byte("$2a$10$Ohlh0UULoGrNO4BAmZzgousbufMk65z9r2hu5Hgbpa8nrAiMQued6")},
		Role:      Role{Name: "admin", Level: 3, Description: "An admin can update and delete other users posts"},
		IsActive:  true,
		CreatedAt: "2025-01-12 17:29:40.000 +0800",
	}, nil
}

func (ms *mockUserStore) Update(context.Context, *User) error {
	return nil
}

func (ms *mockUserStore) DeleteUser(context.Context, int64) error {
	return nil
}

func (ms *mockUserStore) CreateAndInvite(ctx context.Context, user *User, token string, exp time.Duration) error {
	// return nil for success create user, and create invitation records
	if user.Email == "hwteh1997@gmail.com" {
		return common.ErrEmailAlreadyExists
	} else if user.Username == "MehNohNah" {
		return common.ErrUserAlreadyExists
	}
	return nil
}

func (ms *mockUserStore) ActivateUser(context.Context, string) error {
	return nil
}

func (ms *mockUserStore) Delete(context.Context, int64) error {
	return nil
}
