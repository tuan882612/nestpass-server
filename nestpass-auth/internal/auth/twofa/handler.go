package twofa

import (
	"net/http"
	"strconv"
	"time"

	"github.com/tuan882612/apiutils"

	"project/internal/auth"
	"project/internal/auth/email"
	"project/pkg/helpers"
)

// struct for handling two factor authentication requests
type Handler struct {
	twofaService *Service
	prodEnv      bool
}

// NewHandler returns a new handler for two factor authentication requests
func NewHandler(deps *auth.Dependencies) *Handler {
	return &Handler{twofaService: NewService(deps), prodEnv: deps.ProdEnv}
}

// Handles the resend code request
func (h *Handler) ResendCode(w http.ResponseWriter, r *http.Request) {
	input := &auth.Resend{}
	if err := input.Deserialize(r.Body); err != nil {
		apiutils.HandleHttpErrors(w, err)
		return
	}

	if err := h.twofaService.ResendCode(r.Context(), input.Email); err != nil {
		apiutils.HandleHttpErrors(w, err)
		return
	}

	resp := apiutils.NewRes(http.StatusOK, "", nil)
	resp.SendRes(w)
}

// Handles the auth code verifcation
func (h *Handler) Verify(w http.ResponseWriter, r *http.Request) {
	Token := &email.Token{}
	if err := Token.Deserialize(r.Body); err != nil {
		apiutils.HandleHttpErrors(w, err)
		return
	}

	userID, err := helpers.GetUidHeader(r)
	if err != nil {
		apiutils.HandleHttpErrors(w, err)
		return
	}

	mode := r.Header.Get("X-Mode")
	if mode == "" {
		resp := apiutils.NewRes(http.StatusBadRequest, "missing mode header", nil)
		resp.SendRes(w)
		return
	}

	data, retryN, err := h.twofaService.VerifyAuthToken(r.Context(), userID, Token.Token, mode)
	if err != nil {
		w.Header().Set("X-Retry-N", strconv.Itoa(retryN))
		apiutils.HandleHttpErrors(w, err)
		return
	}

	resp := apiutils.NewRes(http.StatusOK, "", nil)
	if mode == "reset" {
		resp.AddHeader(w, map[string]string{"X-Uid": data})
	} else {
		token, err := auth.GenerateStateToken()
		if err != nil {
			apiutils.HandleHttpErrors(w, err)
			return
		}
		// Set JWT as HttpOnly cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "Authorization",
			Value:    data,
			Expires:  time.Now().Add(12 * time.Hour),
			Path:     "/",
			HttpOnly: false,
			Secure:   h.prodEnv,
			SameSite: http.SameSiteNoneMode,
		})
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    token,
			Path:     "/",
			HttpOnly: false,
			Secure:   h.prodEnv,
			SameSite: http.SameSiteNoneMode,
		})
	}
	resp.SendRes(w)
}

// Handles the initial login phase (sending the verification code)
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	input := &auth.Login{}
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

// Handles the initial register phase (sending the verification code)
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	input, err := auth.NewRegister(r.Body)
	if err != nil {
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

// Handles the first reset password phase (sending the verification code)
func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	email := &auth.Resend{}
	if err := email.Deserialize(r.Body); err != nil {
		apiutils.HandleHttpErrors(w, err)
		return
	}

	userID, err := h.twofaService.ResetPassword(r.Context(), email.Email)
	if err != nil {
		apiutils.HandleHttpErrors(w, err)
		return
	}

	resp := apiutils.NewRes(http.StatusOK, "", nil)
	resp.AddHeader(w, map[string]string{"X-Uid": userID})
	resp.SendRes(w)
}

// Handles the third reset password phase (final)
func (h *Handler) ResetPasswordFinal(w http.ResponseWriter, r *http.Request) {
	input := &auth.ResetPsw{}
	if err := input.Deserialize(r.Body); err != nil {
		apiutils.HandleHttpErrors(w, err)
		return
	}

	userID, err := helpers.GetUidHeader(r)
	if err != nil {
		apiutils.HandleHttpErrors(w, err)
		return
	}

	if err := h.twofaService.ResetPasswordFinal(r.Context(), userID, input.Password); err != nil {
		apiutils.HandleHttpErrors(w, err)
		return
	}

	resp := apiutils.NewRes(http.StatusOK, "", nil)
	resp.SendRes(w)
}
