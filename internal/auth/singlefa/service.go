package singlefa

import (
	"context"

	"project/internal/auth"
	"project/internal/auth/jwt"
)

// This is the single-factor authentication service.
type service struct {
	authService auth.Service // base authentication service
	jwtManager  *jwt.Manager
}

// Constructor for the single-factor authentication service.
func NewService(deps *auth.Dependencies) Service {
	return &service{
		authService: deps.Service,
		jwtManager:  deps.JWTManager,
	}
}

// Verifies the user's credentials and returns a JWT token if the verification is successful.
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

// Registers a new user and returns a JWT token if the registration is successful.
func (s *service) SfaRegister(ctx context.Context, input *auth.RegisterInput) (string, error) {
	// convert RegisterInput to RegisterResp and validate the input
	regResp, err := auth.NewRegisterResp(input)
	if err != nil {
		return "", err
	}

	if err := s.authService.RegisterUser(ctx, regResp); err != nil {
		return "", err
	}

	token, err := s.jwtManager.GenerateToken(regResp.UserID)
	if err != nil {
		return "", err
	}

	return token, nil
}
