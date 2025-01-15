package middlewares

import (
	"github.com/juxue97/auth/config"
	"github.com/juxue97/auth/internal/authenticator"
	"github.com/juxue97/auth/internal/cache"
	"github.com/juxue97/auth/internal/repository"
)

type MiddlewareService struct {
	cfg           *config.Config
	PgStore       *repository.Repository
	authenticator *authenticator.Authenticator
	cacheStorage  *cache.RedisCacheStorage
	rateLimiter   *cache.RedisRateLimitStorage
}

func NewMiddlewareService(cfg *config.Config, pgStore *repository.Repository, authenticator *authenticator.Authenticator, cacheStorage *cache.RedisCacheStorage, rateLimiter *cache.RedisRateLimitStorage) *MiddlewareService {
	return &MiddlewareService{
		cfg:           cfg,
		PgStore:       pgStore,
		authenticator: authenticator,
		cacheStorage:  cacheStorage,
		rateLimiter:   rateLimiter,
	}
}
