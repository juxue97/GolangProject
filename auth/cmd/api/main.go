package main

import (
	"github.com/go-redis/redis/v8"
	"github.com/juxue97/auth/config"
	"github.com/juxue97/auth/internal/authenticator"
	"github.com/juxue97/auth/internal/cache"
	"github.com/juxue97/auth/internal/db"
	"github.com/juxue97/auth/internal/mailer"
	"github.com/juxue97/auth/internal/repository"
	"github.com/juxue97/common"
)

//	@title			Auth API
//	@description	This is a authentication backend server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

// @license.name				Apache 2.0
// @license.url				http://www.apache.org/licenses/LICENSE-2.0.html
//
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {
	// Consume .env here

	cfg := config.GetConfig(".env")

	// logger
	logger := common.NewLogger(cfg.Env)
	defer logger.Sync()

	// Pgdb
	pgDB, err := db.NewPgClient(cfg.DB.Addr, cfg.DB.MaxOpenConns, cfg.DB.MaxIdleConns, cfg.DB.MaxIdleTime)
	if err != nil {
		logger.Fatal(err)
	}
	defer pgDB.Close()
	store := repository.NewRepository(pgDB)
	logger.Info("PgStore initialized")

	// Cache
	var redisClient, rateLimitClient *redis.Client
	var cacheStorage *cache.RedisCacheStorage
	var rateLimitStorage *cache.RedisRateLimitStorage

	if cfg.RedisCfg.Enabled {
		redisClient = cache.NewRedisClient(
			cfg.RedisCfg.Addr,
			cfg.RedisCfg.Password,
			cfg.RedisCfg.DB,
		)
		defer redisClient.Close()
		logger.Info("Redis Client initialized")

		if cfg.RateLimit.Enabled {
			rateLimitClient = cache.NewRedisClient(
				cfg.RedisCfg.Addr,
				cfg.RedisCfg.Password,
				cfg.RateLimit.DB,
			)
			defer rateLimitClient.Close()
			logger.Info("Rate Limiter Client initialized")
		}

		cacheStorage = cache.NewRedisStorage(redisClient, cfg.RedisCfg.TTL)
		if cfg.RateLimit.Enabled {
			rateLimitStorage = cache.NewRedisRateLimiterStorage(rateLimitClient)
		}
	}

	// mailtrap
	mailTrapMailer, err := mailer.NewMailTrapClient(cfg.Mail.MailTrap.ApiKey, cfg.Mail.FromEmail)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("Mailer Service initialized")

	// jwtAuthenticator
	jwtAuthenticator := authenticator.NewJwtAuth(
		cfg.Auth.Token.Exp,
		cfg.Auth.Token.Secret,
		cfg.Auth.Token.Aud,
		cfg.Auth.Token.Iss,
	)
	logger.Info("JwtAuthenticator initialized")

	app := &application{
		config:        cfg,
		store:         store,
		cacheStorage:  cacheStorage,
		logger:        logger,
		mailer:        mailTrapMailer,
		authenticator: jwtAuthenticator,
		rateLimiter:   rateLimitStorage,
	}
	mux := app.mount()

	logger.Fatal(app.run(mux))
}
