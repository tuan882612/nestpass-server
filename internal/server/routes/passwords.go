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
		
		r.Route("/password", func(r chi.Router) {
			r.Get("/", handler.GetPassword)
			r.Post("/", handler.CreatePassword)
			r.Patch("/", handler.UpdatePassword)
			r.Delete("/", handler.DeletePassword)
		})
	}
}
