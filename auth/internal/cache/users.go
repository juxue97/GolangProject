package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/juxue97/auth/internal/repository"
)

type UserStore struct {
	rdb *redis.Client
	ttl time.Duration
}

func (us *UserStore) Get(ctx context.Context, userID int64) (*repository.User, error) {
	cacheKey := fmt.Sprintf("user-%d", userID)

	data, err := us.rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var user repository.User
	if data != "" {
		if err := json.Unmarshal([]byte(data), &user); err != nil {
			return nil, err
		}
	}
	return &user, nil
}

func (us *UserStore) Set(ctx context.Context, user *repository.User) error {
	cacheKey := fmt.Sprintf("user-%d", user.ID)
	json, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return us.rdb.SetEX(ctx, cacheKey, json, us.ttl).Err()
}

func (us *UserStore) Delete(ctx context.Context, userID int64) error {
	cacheKey := fmt.Sprintf("user-%d", userID)

	err := us.rdb.Del(ctx, cacheKey).Err()
	if err != nil {
		return fmt.Errorf("failed to delete user from cache: %d", err)
	}

	return nil
}
