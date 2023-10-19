package middlewares

import (
	"context"
	"net/http"

	"project/internal/config"
	"project/pkg/auth"

	"github.com/tuan882612/apiutils"
)

func Authorization(cfg *config.Configuration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := auth.GetBearerToken(r)
			if err != nil {
				apiutils.HandleHttpErrors(w, err)
				return
			}

			claims, err := auth.DecodeToken(token, cfg.SignKey)
			if err != nil {
				apiutils.HandleHttpErrors(w, err)
				return
			}

			// parse and set user id
			ctx := context.WithValue(r.Context(), auth.CtxUserID, claims.UserID)

			// call next handler
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
