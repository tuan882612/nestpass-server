package middlewares

import "github.com/go-chi/cors"

func CorsConfig() cors.Options {
	return cors.Options{
		AllowedOrigins:   []string{
			"http://localhost:5173",
			"http://localhost:5000",
			"http://localhost:3000",
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
	}
}