package oauth

import (
	"net/http"
	"project/internal/config"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/tuan882612/apiutils"
)

type Handler struct {
	svc *Service
}

func NewHandler(cfg *config.Configuration) *Handler {
	return &Handler{svc: NewService(cfg.OAuth)}
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
	c, err := r.Cookie("oauth_state")
	if err != nil {
		resp := apiutils.NewRes(http.StatusBadRequest, "Missing state cookie", nil)
		resp.SendRes(w)
		return
	}

	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	token, err := h.svc.CallbackOAuth(r.Context(), code, state, c.Value)
	if err != nil {
		resp := apiutils.NewRes(http.StatusInternalServerError, "Failed to get token", err.Error())
		resp.SendRes(w)
		return
	}

	resp := apiutils.NewRes(http.StatusOK, "Authenticated successfully", token)
	resp.SendRes(w)
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement token refresh logic here
	resp := apiutils.NewRes(http.StatusOK, "", nil)
	resp.SendRes(w)
}
