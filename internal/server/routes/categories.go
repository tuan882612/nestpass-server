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
		
		r.Route("/category", func(r chi.Router) {
			r.Get("/", handler.GetCategory)
			r.Post("/", handler.CreateCategory)
			r.Patch("/", handler.UpdateCategory)
			r.Delete("/", handler.DeleteCategory)
		})
		r.Route("/passwords", Passwords(deps))
	}
}
