package cache

import (
	"github.com/go-redis/redis/v8"
	"github.com/juxue97/auth/internal/config"
	"github.com/juxue97/common"
)

func NewRedisClient(addr, password string, db int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
}

var (
	RedisClient      *redis.Client
	RateLimitClient  *redis.Client
	CacheStorage     Storage
	RateLimitStorage Storage
)

func init() {
	if config.Configs.RedisCfg.Enabled {
		RedisClient = NewRedisClient(
			config.Configs.RedisCfg.Addr,
			config.Configs.RedisCfg.Password,
			config.Configs.RedisCfg.DB,
		)

		CacheStorage = NewRedisStorage(RedisClient)
		common.Logger.Info("Redis store initialized")

		RateLimitClient = NewRedisClient(
			config.Configs.RedisCfg.Addr,
			config.Configs.RedisCfg.Password,
			config.Configs.RateLimit.DB,
		)

		RateLimitStorage = NewRedisStorage(RateLimitClient)
		common.Logger.Info("Redis rate limit store initialized")
	}
}
