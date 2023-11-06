package routes

import (
	"github.com/go-chi/chi/v5"

	"project/internal/auth/cli"
)

func Cli(handler *cli.Handler, r chi.Router) func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/verify", handler.VerifyCliKey)
		r.Post("/login", handler.Login)
	}
}
