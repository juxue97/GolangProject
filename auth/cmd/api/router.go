package main

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
	"github.com/juxue97/auth/internal/config"
	"github.com/juxue97/common"

	"github.com/juxue97/auth/docs"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type application struct {
	config config.Config
	// store  repository.Repository
	// logger *zap.SugaredLogger
}

const basePath = "/v1"

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

	// apiVersion := fmt.Sprintf("/%s", version)
	r.Route(basePath, func(r chi.Router) {
		if config.Configs.RateLimit.Enabled {
			r.Use(middlewares.RateLimiterMiddleware)
		}
		r.Get("/health", app.healthCheckHandler)

		docsURL := fmt.Sprintf("%s/swagger/doc.json", app.config.Addr)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))

		// Auth routes
		// Public apis
		r.Route("/auth", func(r chi.Router) {
			// Register a new user
			r.Post("/user", auth.RegisterUserHandler)
			// Login user, it can be better if the token stored on the cookies
			r.Post("/login", auth.LoginUserHandler)
		})

		// Users routes
		// Private apis
		r.Route("/users", func(r chi.Router) {
			// Activate the user account
			r.Put("/activate/{token}", users.ActivateUserHandler)
			r.Group(func(r chi.Router) {
				r.Use(middlewares.AuthTokenMiddleware)
				// Get all users
				r.Get("/", middlewares.RoleMiddleware("admin", users.GetUsersHandler))

				r.Route("/{id}", func(r chi.Router) {
					r.Use(middlewares.UsersContextMiddleware)
					r.Get("/", middlewares.RoleMiddleware("admin", users.GetUserHandler))
					r.Put("/", middlewares.RoleMiddleware("admin", users.UpdateUserHandler))
					r.Delete("/", middlewares.RoleMiddleware("admin", users.DeleteUserHandler))
				})
			})

			// Get all users
			// r.Get("/", users.GetUsersHandler)
			// Get a user
			// r.Get("/{id}", users.GetUserHandler)
			// Update a user
			// r.Put("/{id}", users.UpdateUserHandler)
			// Delete a user
			// r.Delete("/{id}", users.DeleteUserHandler)
		})
	})
	// chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
	// 	fmt.Printf("Registered route: %s %s\n", method, route)
	// 	return nil
	// })
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

	common.Logger.Infow("server has started", "addr", app.config.Addr, "enviroment", app.config.Env)

	return srv.ListenAndServe()
}
