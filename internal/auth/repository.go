package auth

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
	"github.com/tuan882612/apiutils"

	"project/internal/config"
	"project/internal/database"
)

type repository struct {
	db *database.DataAccess
}

func NewRepository(cfg *config.Configuration) (Repository, error) {
	if cfg == nil {
		msg := "nil configuration"
		log.Error().Str("location", "NewRepository").Msg(msg)
		return nil, errors.New(msg)
	}

	databases, err := database.NewDataAccess(cfg.NumCpu, cfg.PgUrl, cfg.RedisUrl, cfg.RedisPsw)
	if err != nil {
		return nil, err
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
	row := r.db.Postgres.QueryRow(ctx, UserCredsQuery, email)
	if err := row.Scan(&userID, &password); err != nil {
		if err == pgx.ErrNoRows {
			return uuid.Nil, "", apiutils.NewErrNotFound("user not found")
		}

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

	row := tx.QueryRow(ctx, AddUserQuery,
		input.UserID,
		input.Email,
		input.Name,
		input.Password,
		input.Registered,
		input.UserStatus,
	)

	if err := row.Scan(&input.UserID); err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return apiutils.NewErrConflict("user already exists")
		}
		log.Error().Str("location",	 "AddUser").Msg(err.Error())
		return err
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
