package middlewares

import (
	"context"
	"net/http"

	"github.com/tuan882612/apiutils"

	"nestpass/internal/config"
	"nestpass/pkg/auth"
)

func Authorization(cfg *config.Configuration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// gets the JWT from the authorization header
			token, err := auth.GetBearerToken(r)
			if err != nil {
				apiutils.HandleHttpErrors(w, err)
				return
			}

			// double submit cookie verification
			cookie, err := r.Cookie("session")
			if err != nil {
				apiutils.HandleHttpErrors(w, apiutils.NewErrUnauthorized("Missing JWT cookie"))
				return
			}

			if token != cookie.Value {
				apiutils.HandleHttpErrors(w, apiutils.NewErrUnauthorized("CSRF tokens do not match"))
				return
			}

			// decode JWT token
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
