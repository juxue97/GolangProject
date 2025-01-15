package test_utils

import (
	"testing"

	"github.com/juxue97/auth/config"
	"github.com/juxue97/auth/internal/authenticator"
	"github.com/juxue97/auth/internal/cache"
	"github.com/juxue97/auth/internal/mailer"
	"github.com/juxue97/auth/internal/repository"
	"github.com/juxue97/common"
)

func NewTestApplication(t *testing.T) *Application {
	cfg := config.GetConfig("../../../../.env")
	logger := common.NewLogger(cfg.Env)
	store := repository.NewMockStore()
	cacheStore := cache.NewMockStore()
	mailTrapMailer := mailer.NewMockMailerClient()
	jwtAuthenticator := authenticator.NewMockAuthenticator()
	rateLimitStorage := cache.NewRateLimitMockStore()
	app := &Application{
		Config_test:        cfg,
		Store_test:         store,
		CacheStorage_test:  cacheStore,
		Logger_test:        logger,
		Mailer_test:        mailTrapMailer,
		Authenticator_test: jwtAuthenticator,
		RateLimiter_test:   rateLimitStorage,
	}
	return app
}
