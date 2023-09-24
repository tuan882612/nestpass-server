package singlefa

import (
	"net/http"

	"github.com/tuan882612/apiutils"

	"project/internal/auth"
)

// Single-factor authentication endpoint handlers.
type Handler struct {
	sfaService Service
}

// Constructor for the single-factor authentication endpoint handlers.
func NewHandler(deps *auth.Dependencies) *Handler {
	return &Handler{sfaService: NewService(deps)}
}

// Handles the login request.
// handles http error: 200, 400, 401, 500
func (s *Handler) Login(w http.ResponseWriter, r *http.Request) {
	input := &auth.LoginInput{}
	if err := input.Deserialize(r.Body); err != nil {
		apiutils.HandleHttpErrors(w, err)
		return
	}

	token, err := s.sfaService.SfaLogin(r.Context(), input)
	if err != nil {
		apiutils.HandleHttpErrors(w, err)
		return
	}

	resp := apiutils.NewRes(http.StatusOK, "", nil)
	resp.AddHeader(w, map[string]string{"Authorization": "Bearer " + token})
	resp.SendRes(w)
}

// Handles the register request.
// handles http error: 201, 400, 409, 500
func (s *Handler) Register(w http.ResponseWriter, r *http.Request) {
	input := &auth.RegisterInput{}
	if err := input.Deserialize(r.Body); err != nil {
		apiutils.HandleHttpErrors(w, err)
		return
	}

	token, err := s.sfaService.SfaRegister(r.Context(), input)
	if err != nil {
		apiutils.HandleHttpErrors(w, err)
		return
	}

	resp := apiutils.NewRes(http.StatusCreated, "", nil)
	resp.AddHeader(w, map[string]string{"Authorization": "Bearer " + token})
	resp.SendRes(w)
}
