package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/juxue97/auth/internal/repository"
)

type Storage struct {
	Users interface {
		Get(context.Context, int64) (*repository.User, error)
		Set(context.Context, *repository.User) error
		Delete(context.Context, int64)
	}
	RateLimiter interface {
		Count(context.Context, string, int, time.Duration) (bool, error)
		GetRemainTime(context.Context, string) (time.Duration, error)
	}
}

func NewRedisStorage(rdb *redis.Client) Storage {
	return Storage{
		Users:       &UserStore{rdb: rdb},
		RateLimiter: &RateLimitStore{rdb: rdb},
	}
}
