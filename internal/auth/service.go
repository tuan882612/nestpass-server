package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/tuan882612/apiutils"
	"github.com/tuan882612/apiutils/securityutils"
)

type service struct {
	repository Repository
}

func NewService(repo Repository) (Service, error) {
	if repo == nil {
		msg := "nil repository"
		log.Error().Str("location", "NewService").Msg(msg)
		return nil, errors.New(msg)
	}

	return &service{repository: repo}, nil
}

func (s *service) VerifyUser(ctx context.Context, email, password string) (uuid.UUID, error) {
	if email == "" || password == "" {
		msg := "empty email or password"
		log.Error().Str("location", "VerifyUser").Msg(msg)
		return uuid.Nil, errors.New(msg)
	}

	userID, userPassword, err := s.repository.GetUserCredentials(ctx, email)
	if err != nil {
		return uuid.Nil, err
	}

	if err := securityutils.ValidatePassword(userPassword, password); err != nil {
		return uuid.Nil, apiutils.NewErrUnauthorized(err.Error())
	}

	return userID, nil
}

func (s *service) RegisterUser(ctx context.Context, input *RegisterInput) (uuid.UUID, error) {
	if input == nil {
		msg := "nil input"
		log.Error().Str("location", "RegisterUser").Msg(msg)
		return uuid.Nil, errors.New(msg)
	}

	regResp, err := NewRegisterResp(input)
	if err != nil {
		return uuid.Nil, err
	}

	tx, err := s.repository.startTx(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	defer tx.Rollback(ctx)

	if err := s.repository.AddUser(ctx, tx, regResp); err != nil {
		return uuid.Nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Info().Str("location", "RegisterUser").Msg("failed to commit transaction")
		return uuid.Nil, err
	}

	return regResp.UserID, nil
}
