package auth

import "github.com/go-chi/chi/v5"

func (a *authHandler) RegisterAuthRoutes(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/user", a.registerUserHandler)
		r.Post("/login", a.loginUserHandler)
	})
}
