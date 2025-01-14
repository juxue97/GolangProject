package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"

	"github.com/juxue97/auth/cmd/api/auth"
	"github.com/juxue97/auth/cmd/api/users"
	middlewares "github.com/juxue97/auth/cmd/middleware"
	"github.com/juxue97/auth/config"
	"github.com/juxue97/auth/internal/authenticator"
	"github.com/juxue97/auth/internal/cache"
	"github.com/juxue97/auth/internal/mailer"
	"github.com/juxue97/auth/internal/repository"
	"github.com/juxue97/common"

	"github.com/juxue97/auth/docs"
)

type application struct {
	config        *config.Config
	store         *repository.Repository
	cacheStorage  *cache.RedisCacheStorage
	logger        *zap.SugaredLogger
	mailer        mailer.MailTrapClient
	authenticator *authenticator.JwtAuth
	rateLimiter   *cache.RedisRateLimitStorage
}

const basePath = "/v1"

// var bearerMiddlewareService *middlewares.MiddlewareService

func (app *application) mount() http.Handler {
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
	middlewareService := middlewares.NewMiddlewareService(app.config, app.store, app.authenticator, app.cacheStorage, app.rateLimiter) // use it later

	if app.config.RateLimit.Enabled {
		r.Use(middlewareService.RateLimiterMiddleware)
	}

	// Inject dependencies into handlers
	auth := auth.NewAuthHandler(app.config, app.logger, app.store, app.authenticator, app.mailer)
	user := users.NewUserHandler(middlewareService, app.store, app.cacheStorage)

	r.Route(basePath, func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)
		// Docs Swagger
		docsURL := fmt.Sprintf("%s/swagger/doc.json", app.config.Addr)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))

		// Public apis:
		// Register user, Login User
		auth.RegisterAuthRoutes(r)

		// Private apis:
		// Get activate user, all users, Get user, Update user, Delete user
		user.RegisterUserRoutes(r)
	})
	chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf("Registered route: %s %s\n", method, route)
		return nil
	})
	return r
}

func (app *application) run(mux http.Handler) error {
	// Docs
	docs.SwaggerInfo.Version = app.config.Version
	docs.SwaggerInfo.Host = app.config.ApiUrl
	docs.SwaggerInfo.BasePath = basePath

	srv := &http.Server{
		Addr:         app.config.Addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	app.logger.Infow("server has started", "addr", app.config.Addr, "enviroment", app.config.Env)

	return srv.ListenAndServe()
}
