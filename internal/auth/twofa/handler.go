package twofa

import (
	"net/http"

	"github.com/tuan882612/apiutils"

	"project/internal/auth"
	"project/internal/auth/email"
	"project/pkg/helpers"
)

// struct for handling two factor authentication requests
type Handler struct {
	twofaService Service
}

// NewHandler returns a new handler for two factor authentication requests
func NewHandler(deps *auth.Dependencies) *Handler {
	return &Handler{twofaService: NewService(deps)}
}

// Handles the resend code request
func (h *Handler) ResendCode(w http.ResponseWriter, r *http.Request) {
	input := &email.ResendInput{}
	if err := input.Deserialize(r.Body); err != nil {
		apiutils.HandleHttpErrors(w, err)
		return
	}

	if err := h.twofaService.SendVerificationEmail(r.Context(), input.UserID, input.Email); err != nil {
		apiutils.HandleHttpErrors(w, err)
		return
	}

	resp := apiutils.NewRes(http.StatusOK, "", nil)
	resp.SendRes(w)
}

// Handles the initial login phase (sending the verification code)
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	input := &auth.LoginInput{}
	if err := input.Deserialize(r.Body); err != nil {
		apiutils.HandleHttpErrors(w, err)
		return
	}

	userIDStr, err := h.twofaService.LoginSend(r.Context(), input)
	if err != nil {
		apiutils.HandleHttpErrors(w, err)
		return
	}

	resp := apiutils.NewRes(http.StatusOK, "", nil)
	resp.AddHeader(w, map[string]string{"X-Uid": userIDStr})
	resp.SendRes(w)
}

// Handles the second login phase (verification)
func (h *Handler) LoginVerify(w http.ResponseWriter, r *http.Request) {
	tokenInput := &email.TokenInput{}
	if err := tokenInput.Deserialize(r.Body); err != nil {
		apiutils.HandleHttpErrors(w, err)
		return
	}

	userID, err := helpers.GetUidHeader(r)
	if err != nil {
		apiutils.HandleHttpErrors(w, err)
		return
	}

	authToken, err := h.twofaService.VerifyAuthToken(r.Context(), userID, tokenInput)
	if err != nil {
		apiutils.HandleHttpErrors(w, err)
		return
	}

	resp := apiutils.NewRes(http.StatusOK, "", nil)
	resp.AddHeader(w, map[string]string{"Authorization": "Bearer " + authToken})
	resp.SendRes(w)
}

// Handles the initial register phase (sending the verification code)
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	input := &auth.RegisterInput{}
	if err := input.Deserialize(r.Body); err != nil {
		apiutils.HandleHttpErrors(w, err)
		return
	}

	userID, err := h.twofaService.RegisterSend(r.Context(), input)
	if err != nil {
		apiutils.HandleHttpErrors(w, err)
		return
	}

	resp := apiutils.NewRes(http.StatusOK, "", nil)
	resp.AddHeader(w, map[string]string{"X-Uid": userID})
	resp.SendRes(w)
}

// Handles the second register phase (verification)
func (h *Handler) RegisterVerify(w http.ResponseWriter, r *http.Request) {
	tokenInput := &email.TokenInput{}
	if err := tokenInput.Deserialize(r.Body); err != nil {
		apiutils.HandleHttpErrors(w, err)
		return
	}

	userID, err := helpers.GetUidHeader(r)
	if err != nil {
		apiutils.HandleHttpErrors(w, err)
		return
	}

	authToken, err := h.twofaService.RegisterVerify(r.Context(), userID, tokenInput)
	if err != nil {
		apiutils.HandleHttpErrors(w, err)
		return
	}

	resp := apiutils.NewRes(http.StatusCreated, "", nil)
	resp.AddHeader(w, map[string]string{"Authorization": "Bearer " + authToken})
	resp.SendRes(w)
}
