package routes

import (
	"github.com/go-chi/chi/v5"

	"nestpass/internal/config"
	"nestpass/internal/server/middlewares"
)

func Users(handler *APIHandler, cfg *config.Configuration) func(r chi.Router) {
	return func(r chi.Router) {

		r.Use(middlewares.Authorization(cfg))
		r.Get("/", handler.User.GetUser)
		r.Get("/clikey", handler.User.GetCliKey)
		r.Put("/clikey", handler.User.CreateCliKey)
		r.Route("/categories", Categories(handler))
	}
}
