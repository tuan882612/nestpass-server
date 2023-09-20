package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/tuan882612/apiutils"
	"github.com/tuan882612/apiutils/securityutils"
)

// This defines the interface for the base authentication service used by other authentication services.
type service struct {
	repository Repository
}

// This is the constructor for the base authentication service.
func NewService(repo Repository) (Service, error) {
	if repo == nil {
		msg := "nil repository"
		log.Error().Str("location", "NewService").Msg(msg)
		return nil, errors.New(msg)
	}

	return &service{repository: repo}, nil
}

// This method retrieves the user credentials and returns the user ID if the credentials are valid.
func (s *service) VerifyUser(ctx context.Context, email, password string) (uuid.UUID, error) {
	if email == "" || password == "" {
		msg := "empty email or password"
		log.Error().Str("location", "VerifyUser").Msg(msg)
		return uuid.Nil, errors.New(msg)
	}

	// retrieve the user credentials from the database
	userID, userPassword, err := s.repository.GetUserCredentials(ctx, email)
	if err != nil {
		return uuid.Nil, err
	}

	// validate the password
	if err := securityutils.ValidatePassword(userPassword, password); err != nil {
		return uuid.Nil, apiutils.NewErrUnauthorized(err.Error())
	}

	return userID, nil
}

// This method registers a new user and returns the user ID if the registration is successful.
func (s *service) RegisterUser(ctx context.Context, input *RegisterInput) (uuid.UUID, error) {
	if input == nil {
		msg := "nil input"
		log.Error().Str("location", "RegisterUser").Msg(msg)
		return uuid.Nil, errors.New(msg)
	}

	// convert RegisterInput to RegisterResp and validate the input
	regResp, err := NewRegisterResp(input)
	if err != nil {
		return uuid.Nil, err
	}

	// start a new transaction and defer rollback if there is an error
	tx, err := s.repository.startTx(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	defer tx.Rollback(ctx)

	// tries to add the user to the database
	if err := s.repository.AddUser(ctx, tx, regResp); err != nil {
		return uuid.Nil, err
	}

	// commit the transaction
	if err := tx.Commit(ctx); err != nil {
		log.Info().Str("location", "RegisterUser").Msg("failed to commit transaction")
		return uuid.Nil, err
	}

	return regResp.UserID, nil
}
