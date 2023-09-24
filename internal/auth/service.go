package auth

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/tuan882612/apiutils"
	"github.com/tuan882612/apiutils/securityutils"
)

// Base authentication service.
type service struct {
	repository Repository
}

// Constructor for the base authentication service.
func NewService(repo Repository) Service {
	return &service{repository: repo}
}

// Retrieves the user credentials and returns the user ID if the credentials are valid.
func (s *service) VerifyUser(ctx context.Context, email, password string) (uuid.UUID, error) {
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

// Registers a new user and returns the user ID if the registration is successful.
func (s *service) RegisterUser(ctx context.Context, body *RegisterResp) error {
	// start a new transaction and rollback if there is an error
	tx, err := s.repository.startTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// tries to add the user to the database
	if err := s.repository.AddUser(ctx, tx, body); err != nil {
		return err
	}

	// commit the transaction
	if err := tx.Commit(ctx); err != nil {
		log.Info().Str("location", "RegisterUser").Msg("failed to commit transaction")
		return err
	}

	return nil
}
