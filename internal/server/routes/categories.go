package routes

import (
	"github.com/go-chi/chi/v5"
)

func Categories(handler *APIHandler) func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/", handler.Category.GetAllCategories)

		r.Route("/category", func(r chi.Router) {
			r.Get("/", handler.Category.GetCategory)
			r.Post("/", handler.Category.CreateCategory)
			r.Patch("/", handler.Category.UpdateCategory)
			r.Delete("/", handler.Category.DeleteCategory)
		})
		r.Route("/passwords", Passwords(handler))
	}
}
