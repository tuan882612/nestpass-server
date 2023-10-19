package middlewares

import (
	"net/http"

	"project/internal/config"
)

func Authorization(next http.Handler, cfg *config.Configuration) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
