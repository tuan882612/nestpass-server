package singlefa

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/tuan882612/apiutils"

	"project/internal/auth"
)

type service struct {
	authService auth.Service
	jwtHandler  *auth.JWTManger
}

func NewService(authSvc auth.Service, jwtHandler *auth.JWTManger) (Service, error) {
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

func (s *service) SfaLogin(ctx context.Context, input *auth.LoginInput) (string, error) {
	userID, err := s.authService.VerifyUser(ctx, input.Email, input.Password)
	if err != nil {
		return "", err
	}

	token, err := s.jwtHandler.GenerateToken(userID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *service) SfaRegister(ctx context.Context, input *auth.RegisterInput) (string, error) {
	userID, err := s.authService.RegisterUser(ctx, input)
	if err != nil {
		return "", err
	}

	token, err := s.jwtHandler.GenerateToken(userID)
	if err != nil {
		return "", err
	}

	return token, nil
}
