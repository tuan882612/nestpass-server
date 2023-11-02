package middlewares

import (
	"context"
	"net/http"
	"net/url"

	"github.com/tuan882612/apiutils"

	"nestpass/internal/config"
	"nestpass/pkg/auth"
)

func Authorization(cfg *config.Configuration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// double submit cookie verification
			csrfToken, err := r.Cookie("token")
			if err != nil {
				apiutils.HandleHttpErrors(w, apiutils.NewErrUnauthorized("Missing CSRF token"))
				return
			}

			csrfTokenHeader, err := url.QueryUnescape(r.Header.Get("X-CSRF-Token"))
			if err != nil {
				apiutils.HandleHttpErrors(w, apiutils.NewErrUnauthorized("CSRF token is invalid"))
				return
			}

			if csrfTokenHeader != csrfToken.Value {
				apiutils.HandleHttpErrors(w, apiutils.NewErrUnauthorized("CSRF tokens do not match"))
				return
			}

			// gets AuthToken from cookie
			cookie, err := r.Cookie("Authorization")
			if err != nil {
				err := apiutils.NewErrUnauthorized("missing authorization cookie")
				apiutils.HandleHttpErrors(w, err)
				return
			}
			token := cookie.Value

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
