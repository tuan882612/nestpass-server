package routes

import (
	"github.com/go-chi/chi/v5"

	"project/internal/auth/oauth"
)

func OAuth(handler oauth.Handler, r chi.Router) func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/", handler.Invoke)
		r.Get("/callback", handler.Callback)
	}
}
