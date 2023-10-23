package routes

import (
	"github.com/go-chi/chi/v5"

	"nestpass/internal/dependencies"
	"nestpass/internal/users/passwords"
)

func Passwords(deps *dependencies.Dependencies) func(r chi.Router) {
	handler := passwords.NewHandler(deps)

	return func(r chi.Router) {
		r.Get("/", handler.GetAllPasswords)
		r.Post("/", handler.CreatePassword)

		r.Route("/{pid}", func(r chi.Router) {
			r.Get("/", handler.GetPassword)
			r.Put("/", handler.UpdatePassword)
			r.Delete("/", handler.DeletePassword)
		})
	}
}
