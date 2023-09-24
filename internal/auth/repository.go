package auth

import (
	"context"
	"errors"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog/log"
	"github.com/tuan882612/apiutils"

	"project/internal/auth/email"
	"project/internal/config"
	"project/internal/database"
)

// Base authentication repository.
type repository struct {
	db *database.DataAccess
}

// Constructor for the base authentication repository.
func NewRepository(cfg *config.Configuration) (Repository, error) {
	databases, err := database.NewDataAccess(cfg.NumCpu, cfg.PgUrl, cfg.RedisUrl, cfg.RedisPsw)
	if err != nil {
		return nil, err
	}

	return &repository{db: databases}, nil
}

// Retrieves the user's uuid and password from the database if the user exists.
func (r *repository) GetUserCredentials(ctx context.Context, email string) (uuid.UUID, string, error) {
	// initialize credential variables
	var userID uuid.UUID
	var password string
	row := r.db.Postgres.QueryRow(ctx, UserCredsQuery, email)

	// scan the row and check for errors
	if err := row.Scan(&userID, &password); err != nil {
		if err == pgx.ErrNoRows {
			return uuid.Nil, "", apiutils.NewErrNotFound("user not found")
		}

		log.Error().Str("location", "GetUserCredentials").Msg(err.Error())
		return uuid.Nil, "", err
	}

	return userID, password, nil
}

// Retrieves the user's twofa data.
func (r *repository) GetTwofaCache(ctx context.Context, userID uuid.UUID) (*email.TwofaBody, error) {
	data, err := r.db.Redis.Get(userID.String()).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, apiutils.NewErrNotFound("twofa not found")
		}

		log.Error().Str("location", "GetTwofaCache").Msg(err.Error())
		return nil, err
	}

	tfaBody := &email.TwofaBody{}
	if err := tfaBody.Deserialize(data); err != nil {
		log.Error().Str("location", "GetTwofaCache").Msg(err.Error())
		return nil, err
	}

	return tfaBody, nil
}

// Updates the user's twofa data.
func (r *repository) UpdateTwofaCache(ctx context.Context, userID uuid.UUID, body *email.TwofaBody) error {
	data, err := body.Serialize()
	if err != nil {
		log.Error().Str("location", "UpdateTwofaCache").Msg(err.Error())
		return err
	}

	// update the twofa data and set the ttl to the previous value
	idStr := userID.String()
	ttl := r.db.Redis.TTL(idStr).Val()
	if err := r.db.Redis.Set(idStr, data, ttl).Err(); err != nil {
		log.Error().Str("location", "UpdateTwofaCache").Msg(err.Error())
		return err
	}

	return nil
}

// Deletes the user's twofa data.
func (r *repository) DeleteTwofaCache(ctx context.Context, userID uuid.UUID) error {
	if err := r.db.Redis.Del(uuid.UUID.String(userID)).Err(); err != nil {
		log.Error().Str("location", "DeleteTwofaCache").Msg(err.Error())
		return err
	}

	return nil
}

// Adds a new user to the database.
func (r *repository) AddUser(ctx context.Context, tx pgx.Tx, input *RegisterResp) error {
	_, err := tx.Exec(ctx, AddUserQuery,
		&input.UserID,
		&input.Email,
		&input.Name,
		&input.Password,
		&input.Registered,
		&input.UserStatus,
	)

	// checking for errors
	if err != nil {
		// initializing pgx error and checking for duplicate key error
		pgErr := &pgconn.PgError{}
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return apiutils.NewErrConflict("user already exists")
		}

		log.Error().Str("location", "AddUser").Msg(err.Error())
		return err
	}

	return nil
}

// Updates the user's status to "active".
func (r *repository) UpdateUserStatus(ctx context.Context, userID uuid.UUID) error {
	if _, err := r.db.Postgres.Exec(ctx, UpdateUserStatusQuery, userID); err != nil {
		log.Error().Str("location", "UpdateUserStatus").Msg(err.Error())
		return err
	}

	return nil
}

// Starts a new postgres transaction.
func (r *repository) startTx(ctx context.Context) (pgx.Tx, error) {
	tx, err := r.db.Postgres.Begin(ctx)
	if err != nil {
		log.Error().Str("location", "startTx").Msg(err.Error())
		return nil, err
	}

	return tx, nil
}
