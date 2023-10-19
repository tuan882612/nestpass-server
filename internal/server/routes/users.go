package routes

import (
	"net/http"
	"project/internal/config"
	"project/internal/server/middlewares"
	"project/pkg/auth"

	"github.com/go-chi/chi/v5"
	"github.com/tuan882612/apiutils"
)

func Users(cfg *config.Configuration, r chi.Router) func(r chi.Router) {
	return func(r chi.Router) {
		r.Use(middlewares.Authorization(cfg))
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			userID, err := auth.UidFromCtx(r.Context())
			if err != nil {
				apiutils.HandleHttpErrors(w, err)
				return
			}

			resp := apiutils.NewRes(http.StatusOK, "success", userID)
			resp.SendRes(w)
		})
	}
}
