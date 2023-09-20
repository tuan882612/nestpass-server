package singlefa

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/tuan882612/apiutils"

	"project/internal/auth"
	"project/internal/auth/jwt"
)

// This is the single-factor authentication service.
type service struct {
	authService auth.Service
	jwtManager  *jwt.JWTManger
}

// This is the constructor for the single-factor authentication service.
// It takes the base authentication service and the JWT manager as dependencies.
func NewService(authSvc auth.Service, jwtManager *jwt.JWTManger) (Service, error) {
	depMap := apiutils.Dependencies{
		"authService": authSvc,
		"jwtManger":   jwtManager,
	}

	if err := apiutils.ValidateDependencies(depMap); err != nil {
		log.Error().Err(err).Msg("failed to validate dependencies")
		return nil, err
	}

	return &service{
		authService: authSvc,
		jwtManager:  jwtManager,
	}, nil
}

// This method verifies the user's credentials and returns a JWT token if the verification is successful.
func (s *service) SfaLogin(ctx context.Context, input *auth.LoginInput) (string, error) {
	userID, err := s.authService.VerifyUser(ctx, input.Email, input.Password)
	if err != nil {
		return "", err
	}

	token, err := s.jwtManager.GenerateToken(userID)
	if err != nil {
		return "", err
	}

	return token, nil
}

// This method registers a new user and returns a JWT token if the registration is successful.
func (s *service) SfaRegister(ctx context.Context, input *auth.RegisterInput) (string, error) {
	userID, err := s.authService.RegisterUser(ctx, input)
	if err != nil {
		return "", err
	}

	token, err := s.jwtManager.GenerateToken(userID)
	if err != nil {
		return "", err
	}

	return token, nil
}
