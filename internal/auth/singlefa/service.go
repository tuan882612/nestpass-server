package singlefa

import (
	"github.com/rs/zerolog/log"
	"github.com/tuan882612/apiutils"

	"project/internal/auth"
)

type service struct {
	authService *auth.Service
	jwtHandler  *auth.JWTManger
}

func NewService(authSvc *auth.Service, jwtHandler *auth.JWTManger) (Service, error) {
	depMap := apiutils.Dependencies{
		"authService": authSvc,
		"jwtHandler":  jwtHandler,
	}

	if err := apiutils.ValidateDependencies(depMap); err != nil {
		log.Error().Err(err).Msg("failed to validate dependencies")
		return nil, err
	}

	return &service{
		authService: authSvc,
		jwtHandler:  jwtHandler,
	}, nil
}

func (s *service) SfaLogin() {

}

func (s *service) SfaRegister() {

}
