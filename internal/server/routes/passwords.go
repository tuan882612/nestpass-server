package routes

import (
	"github.com/go-chi/chi/v5"
)

func Passwords(handler *APIHandler) func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/", handler.Password.GetAllPasswords)

		r.Route("/password", func(r chi.Router) {
			r.Get("/", handler.Password.GetPassword)
			r.Post("/", handler.Password.CreatePassword)
			r.Patch("/", handler.Password.UpdatePassword)
			r.Delete("/", handler.Password.DeletePassword)
		})
	}
}
