package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"

	"project/internal/database"
)

type repository struct {
	db *database.DataAccess
}

func NewRepository(databases *database.DataAccess) (Repository, error) {
	if databases == nil {
		msg := "nil database"
		log.Error().Str("location", "NewRepository").Msg(msg)
		return nil, errors.New(msg)
	}

	return &repository{db: databases}, nil
}

func (r *repository) GetUserCredentials(ctx context.Context, email string) (uuid.UUID, string, error) {
	if email == "" {
		msg := "empty email"
		log.Error().Str("location", "GetUserCredentials").Msg(msg)
		return uuid.Nil, "", errors.New(msg)
	}

	var userID uuid.UUID
	var password string
	row := r.db.Postgres.QueryRow(ctx, UserCreds)
	if err := row.Scan(&userID, &password); err != nil {
		log.Error().Str("location", "GetUserCredentials").Msg(err.Error())
		return uuid.Nil, "", err
	}

	return userID, password, nil
}

func (r *repository) AddUser(ctx context.Context, tx pgx.Tx, input *RegisterResp) error {
	if input == nil {
		msg := "nil input"
		log.Error().Str("location", "AddUser").Msg(msg)
		return errors.New(msg)
	}

	return nil
}

func (r *repository) startTx(ctx context.Context) (pgx.Tx, error) {
	tx, err := r.db.Postgres.Begin(ctx)
	if err != nil {
		log.Error().Str("location", "startTx").Msg(err.Error())
		return nil, err
	}

	return tx, nil
}
