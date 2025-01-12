package cache

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/juxue97/auth/internal/repository"
)

type Storage struct {
	Users interface {
		Get(context.Context, int64) (*repository.User, error)
		Set(context.Context, *repository.User) error
		Delete(context.Context, int64)
	}
}

func NewRedisStorage(rdb *redis.Client) Storage {
	return Storage{
		Users: &UserStore{rdb: rdb},
	}
}
