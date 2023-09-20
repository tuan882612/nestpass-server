package singlefa

import (
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/tuan882612/apiutils"

	"project/internal/auth"
	"project/internal/auth/jwt"
)

// This is the single-factor authentication endpoint handlers.
type Handler struct {
	sfaService Service
}

// This is the constructor for the single-factor authentication endpoint handlers.
func NewHandler(authService auth.Service, jwtHandler *jwt.JWTManger) (*Handler, error) {
	depMap := apiutils.Dependencies{
		"authService": authService,
		"jwtHandler":  jwtHandler,
	}

	if err := apiutils.ValidateDependencies(depMap); err != nil {
		log.Error().Err(err).Msg("failed to validate dependencies")
		return nil, err
	}

	service, err := NewService(authService, jwtHandler)
	if err != nil {
		return nil, err
	}

	return &Handler{
		sfaService: service,
	}, nil
}

// This method handles the login request.
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

// This method handles the register request.
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
