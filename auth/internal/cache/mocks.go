package cache

import (
	"context"
	"time"

	"github.com/juxue97/auth/internal/repository"
)

type (
	newMockUserStore      struct{}
	newMockRateLimitStore struct{}
)

func NewMockStore() *RedisCacheStorage {
	return &RedisCacheStorage{
		Users: &newMockUserStore{},
	}
}

func NewRateLimitMockStore() *RedisRateLimitStorage {
	return &RedisRateLimitStorage{
		RateLimiter: &newMockRateLimitStore{},
	}
}

func (rls *newMockRateLimitStore) Count(context.Context, string, int, time.Duration) (bool, error) {
	return true, nil
}

func (rls *newMockRateLimitStore) GetRemainTime(context.Context, string) (time.Duration, error) {
	return time.Hour, nil
}

func (ms *newMockUserStore) Get(context.Context, int64) (*repository.User, error) {
	return nil, nil
}

func (ms *newMockUserStore) Set(context.Context, *repository.User) error {
	return nil
}

func (ms *newMockUserStore) Delete(context.Context, int64) error {
	return nil
}
