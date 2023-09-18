package routes

import (
	"github.com/go-chi/chi/v5"

	"project/internal/auth/singlefa"
)

func SingleFA(handler singlefa.Handler) func(r chi.Router) {
	return func(r chi.Router) {
		r.Post("/login", handler.Login)
		r.Post("/register", handler.Register)
	}
}
