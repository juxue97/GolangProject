package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/juxue97/auth/cmd/api/users"
	"github.com/juxue97/auth/internal/db"
	"github.com/juxue97/auth/internal/repository"
	"go.uber.org/zap"
)

type application struct {
	config config
	store  repository.Repository
	logger *zap.SugaredLogger
}

type config struct {
	url     string
	version string
	addr    string
	env     string
	db      db.PgDBConfig
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))
	apiVersion := fmt.Sprintf("/%s", version)

	r.Route(apiVersion, func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)

		r.Route("/users", func(r chi.Router) {
			r.Post("/", users.CreateUserHandler)
		})
	})

	return r
}

func (app *application) run(mux http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}
	log.Printf("server has started at %s", app.config.addr)

	return srv.ListenAndServe()
}
