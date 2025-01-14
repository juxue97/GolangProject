package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/juxue97/auth/internal/repository"
)

type RedisCacheStorage struct {
	Users interface {
		Get(context.Context, int64) (*repository.User, error)
		Set(context.Context, *repository.User) error
		Delete(context.Context, int64) error
	}
}

func NewRedisStorage(rdb *redis.Client, ttl time.Duration) *RedisCacheStorage {
	return &RedisCacheStorage{
		Users: &UserStore{rdb: rdb, ttl: ttl},
	}
}

type RedisRateLimitStorage struct {
	RateLimiter interface {
		Count(context.Context, string, int, time.Duration) (bool, error)
		GetRemainTime(context.Context, string) (time.Duration, error)
	}
}

func NewRedisRateLimiterStorage(rdb *redis.Client) *RedisRateLimitStorage {
	return &RedisRateLimitStorage{
		RateLimiter: &RateLimitStore{rdb: rdb},
	}
}
