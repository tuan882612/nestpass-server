package routes

import (
	"github.com/go-chi/chi/v5"

	"nestpass/internal/config"
	"nestpass/internal/dependencies"
	"nestpass/internal/server/middlewares"
	"nestpass/internal/users"
)

func Users(deps *dependencies.Dependencies, cfg *config.Configuration) func(r chi.Router) {
	userHandler := users.NewHandler(deps)

	return func(r chi.Router) {
		
		r.Use(middlewares.Authorization(cfg))
		r.Get("/", userHandler.GetUser)
		r.Get("/clikey", userHandler.GetCliKey)
		r.Put("/clikey", userHandler.CreateCliKey)
		r.Route("/categories", Categories(deps))
	}
}
