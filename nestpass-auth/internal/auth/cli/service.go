package cli

import (
	"context"

	"github.com/google/uuid"
	"github.com/tuan882612/apiutils"
	"github.com/tuan882612/apiutils/securityutils"

	"project/internal/auth"
	"project/internal/auth/jwt"
)

// Service for handling cli authentication.
type Service struct {
	authRepo   *auth.Repository // base auth repository
	cacheRepo  *auth.Cache      // cache repository
	jwtManager *jwt.Manager
}

// Creates a new cli authentication service with the given dependencies.
func NewService(deps *auth.Dependencies) *Service {
	return &Service{
		authRepo:   deps.Repository,
		cacheRepo:  deps.Cache,
		jwtManager: deps.JWTManager,
	}
}

// Verifies the cli key
func (s *Service) VerifyCliKey(ctx context.Context, userID uuid.UUID, inputCliKey string) (string, error) {
	// retrieve the cli key from the cache and compare it with the input cli key
	cliKey, err := s.cacheRepo.GetData(ctx, userID, auth.Cli)
	if err != nil {
		return "", err
	}

	if cliKey != inputCliKey {
		return "", apiutils.NewErrUnauthorized("invalid clikey")
	}

	// generate a new jwt token
	token, err := s.jwtManager.GenerateToken(userID)
	if err != nil {
		return "", err
	}

	return token, nil
}

// Initial twofa login
func (s *Service) LoginSend(ctx context.Context, input *auth.Login) (string, error) {
	// retrieve the user credentials from the database
	user, err := s.authRepo.GetUserCredentials(ctx, input.Email)
	if err != nil {
		return "", err
	}

	// check if user is oauth user
	if user.Password == "" {
		return "", apiutils.NewErrBadRequest("user is oauth user")
	}

	// validate the password
	if err := securityutils.ValidatePassword(user.Password, input.Password); err != nil {
		return "", apiutils.NewErrUnauthorized(err.Error())
	}

	return user.UserID.String(), nil
}
