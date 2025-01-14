package cache

import (
	"github.com/go-redis/redis/v8"
)

var (
	RedisClient      *redis.Client
	RateLimitClient  *redis.Client
	CacheStorage     RedisCacheStorage
	RateLimitStorage RedisRateLimitStorage
)

func NewRedisClient(addr, password string, db int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
}
