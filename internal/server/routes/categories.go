package routes

import (
	"github.com/go-chi/chi/v5"

	"nestpass/internal/dependencies"
	"nestpass/internal/users/categories"
)

func Categories(deps *dependencies.Dependencies) func(r chi.Router) {
	handler := categories.NewHandler(deps)

	return func(r chi.Router) {
		r.Get("/", handler.GetAllCategories)
		r.Post("/", handler.CreateCategory)

		r.Route("/{cid}", func(r chi.Router) {
			r.Get("/", handler.GetCategory)
			r.Put("/", handler.UpdateCategory)
			r.Delete("/", handler.DeleteCategory)

			r.Route("/passwords", Passwords(deps))
		})
	}
}
