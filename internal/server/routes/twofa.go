package routes

import (
	"github.com/go-chi/chi/v5"

	"project/internal/auth/twofa"
)

func TwoFA(handler *twofa.Handler, r chi.Router) func(r chi.Router) {
	return func(r chi.Router) {
		r.Post("/resend", handler.ResendCode)
		r.Post("/verify", handler.Verify)
		r.Post("/login", handler.Login)
		r.Post("/register", handler.Register)
		r.Post("/reset", handler.ResetPassword)
		r.Post("/reset/verify", handler.ResetPasswordVerify)
		r.Post("/reset/final", handler.ResetPasswordFinal)
	}
}
