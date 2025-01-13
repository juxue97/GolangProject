package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type RateLimitStore struct {
	rdb *redis.Client
}

func (rls *RateLimitStore) Count(ctx context.Context, ip string, limit int, window time.Duration) (bool, error) {
	// Define the Redis key for the IP
	key := fmt.Sprintf("ratelimit:%s", ip)

	// Start a Redis transaction
	pipe := rls.rdb.TxPipeline()

	// Increment the count for the IP
	incr := pipe.Incr(ctx, key)

	ttl := rls.rdb.TTL(ctx, key)
	if ttl.Val() == -2 {
		// Set an expiration on the key only if it is newly created
		pipe.Expire(ctx, key, window)
	}

	// Execute the transaction
	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}

	// Check if the user has exceeded the limit
	if incr.Val() > int64(limit) {
		return false, nil // Rate limit exceeded
	}

	return true, nil // Request allowed
}

func (rls *RateLimitStore) GetRemainTime(ctx context.Context, ip string) (time.Duration, error) {
	key := fmt.Sprintf("ratelimit:%s", ip)

	ttl, err := rls.rdb.TTL(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	return ttl, nil
}
