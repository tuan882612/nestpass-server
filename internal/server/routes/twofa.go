package routes

import (
	"github.com/go-chi/chi/v5"

	"project/internal/auth/twofa"
)

func TwoFA(handler *twofa.Handler, r chi.Router) func(r chi.Router) {
	return func(r chi.Router) {
		r.Post("/resend", handler.ResendCode)
		r.Post("/login", handler.Login)
		r.Post("/login/verify", handler.LoginVerify)
		r.Post("/register", handler.Register)
		r.Post("/register/verify", handler.RegisterVerify)
		r.Post("/resend", handler.ResendCode)
		r.Post("/reset", handler.ResetPassword)
		r.Post("/reset/verify", handler.ResetPasswordVerify)
	}
}
