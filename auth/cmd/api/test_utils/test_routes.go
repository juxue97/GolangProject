package test_utils

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/juxue97/auth/cmd/api/auth"
	"github.com/juxue97/auth/cmd/api/users"
	middlewares "github.com/juxue97/auth/cmd/middleware"
	"github.com/juxue97/auth/config"
	"github.com/juxue97/auth/internal/authenticator"
	"github.com/juxue97/auth/internal/cache"
	"github.com/juxue97/auth/internal/mailer"
	"github.com/juxue97/auth/internal/repository"
	"github.com/juxue97/common"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

type Application struct {
	Config_test        *config.Config
	Store_test         *repository.Repository
	CacheStorage_test  *cache.RedisCacheStorage
	Logger_test        *zap.SugaredLogger
	Mailer_test        *mailer.Client
	Authenticator_test *authenticator.Authenticator
	RateLimiter_test   *cache.RedisRateLimitStorage
}

const basePath = "/v1"

func (app *Application) Mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{common.GetString("CORS_ALLOWED_ORIGIN", "http://localhost:3001")},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.StripSlashes)

	// Inject dependencies into services
	middlewareService := middlewares.NewMiddlewareService(app.Config_test, app.Store_test, app.Authenticator_test, app.CacheStorage_test, app.RateLimiter_test)

	if app.Config_test.RateLimit.Enabled {
		r.Use(middlewareService.RateLimiterMiddleware)
	}

	// Inject dependencies into handlers
	auth := auth.NewAuthHandler(app.Config_test, app.Logger_test, app.Store_test, app.Authenticator_test, app.Mailer_test)
	user := users.NewUserHandler(middlewareService, app.Store_test, app.CacheStorage_test)

	r.Route(basePath, func(r chi.Router) {
		// Docs Swagger
		docsURL := fmt.Sprintf("%s/swagger/doc.json", app.Config_test.Addr)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))

		// Public apis:
		// Register user, Login User
		auth.RegisterAuthRoutes(r)

		// Private apis:
		// Get activate user, all users, Get user, Update user, Delete user
		// Modify here for mock, remove middleware

		user.RegisterUserRoutes(r)
	})
	// chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
	// 	fmt.Printf("Registered route: %s %s\n", method, route)
	// 	return nil
	// })
	return r
}
