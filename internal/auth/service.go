package auth

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
)

type service struct {
	repository Repository
}

func NewService(authRepo Repository) (Service, error) {
	if authRepo == nil {
		msg := "nil repository"
		log.Error().Str("location", "NewService").Msg(msg)
		return nil, errors.New(msg)
	}

	return &service{repository: authRepo}, nil
}

func (s *service) VerifyUser(ctx context.Context, email, password string) (string, error) {
	return "", nil
}

func (s *service) RegisterUser(ctx context.Context, input RegisterInput) (string, error) {
	return "", nil
}
