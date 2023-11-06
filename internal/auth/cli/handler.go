package cli

import (
	"net/http"

	"github.com/tuan882612/apiutils"

	"project/internal/auth"
	"project/pkg/helpers"
)

// struct for handling two factor authentication requests
type Handler struct {
	cliService *Service
}

// NewHandler returns a new handler for two factor authentication requests
func NewHandler(deps *auth.Dependencies) *Handler {
	return &Handler{cliService: NewService(deps)}
}

// Handles verifying the cli key
func (h *Handler) VerifyCliKey(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := helpers.GetUidHeader(r)
	if err != nil {
		apiutils.HandleHttpErrors(w, err)
		return
	}

	cliKey := r.Header.Get("X-CLI-Key")
	if cliKey == "" {
		apiutils.HandleHttpErrors(w, apiutils.NewErrBadRequest("missing clikey header"))
		return
	}

	if err := h.cliService.VerifyCliKey(ctx, userID, cliKey); err != nil {
		apiutils.HandleHttpErrors(w, err)
		return
	}

	resp := apiutils.NewRes(http.StatusOK, "", nil)
	resp.SendRes(w)
}

// Handles the initial login phase (sending the verification code)
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	input := &auth.Login{}
	if err := input.Deserialize(r.Body); err != nil {
		apiutils.HandleHttpErrors(w, err)
		return
	}

	userIDStr, err := h.cliService.LoginSend(r.Context(), input)
	if err != nil {
		apiutils.HandleHttpErrors(w, err)
		return
	}

	resp := apiutils.NewRes(http.StatusOK, "", nil)
	resp.AddHeader(w, map[string]string{"X-Uid": userIDStr})
	resp.SendRes(w)
}
