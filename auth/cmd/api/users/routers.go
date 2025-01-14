package users

import "github.com/go-chi/chi/v5"

func (u *userHandler) RegisterUserRoutes(r chi.Router) {
	r.Route("/users", func(r chi.Router) {
		r.Put("/activate/{token}", u.activateUserHandler)
		r.Group(func(r chi.Router) {
			r.Use(u.middlewareService.AuthTokenMiddleware)
			// Get all users
			r.Get("/", u.middlewareService.RoleMiddleware("admin", u.getUsersHandler))

			r.Route("/{id}", func(r chi.Router) {
				r.Use(u.middlewareService.UsersContextMiddleware)
				r.Get("/", u.middlewareService.RoleMiddleware("admin", u.getUserHandler))
				r.Put("/", u.middlewareService.RoleMiddleware("admin", u.updateUserHandler))
				r.Delete("/", u.middlewareService.RoleMiddleware("admin", u.deleteUserHandler))
			})
		})
	})
}
