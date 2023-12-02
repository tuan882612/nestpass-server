package oauth

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/tuan882612/apiutils"

	"project/internal/auth"
	"project/internal/config"
)

type Handler struct {
	svc *Service
}

func NewHandler(cfg *config.Configuration, deps *auth.Dependencies) *Handler {
	return &Handler{svc: NewService(cfg.OAuth, deps)}
}

func (h *Handler) Invoke(w http.ResponseWriter, r *http.Request) {
	state := h.svc.generateStateOauthCookie()
	http.SetCookie(w, &http.Cookie{
		Name:    "oauth_state",
		Value:   state,
		Expires: time.Now().Add(10 * time.Minute),
	})

	redirectURL := h.svc.StartOAuth(r.Context(), state)
	http.Redirect(w, r, redirectURL, http.StatusFound)
	log.Info().Msg("Redirected to OAuth provider")
}

func (h *Handler) Callback(w http.ResponseWriter, r *http.Request) {
	// get state cookie
	c, err := r.Cookie("oauth_state")
	if err != nil {
		resp := apiutils.NewRes(http.StatusBadRequest, "Missing state cookie", nil)
		resp.SendRes(w)
		return
	}

	// check if state cookie matches the state query param
	ctx := r.Context()
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	token, err := h.svc.CallbackOAuth(ctx, code, state, c.Value)
	if err != nil {
		apiutils.HandleHttpErrors(w, err)
		return
	}

	// invalidate state cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "oauth_state",
		Value:   "",
		Expires: time.Unix(0, 0),
		MaxAge:  -1,
	})

	authToken, err := h.svc.UserLoginSignup(ctx, token)
	if err != nil {
		apiutils.HandleHttpErrors(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    authToken,
		Expires:  time.Now().Add(12 * time.Hour),
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	})

	resp := apiutils.NewRes(http.StatusOK, "authenticated successfully", nil)
	resp.SendRes(w)
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement token refresh logic here
	resp := apiutils.NewRes(http.StatusOK, "", nil)
	resp.SendRes(w)
}
