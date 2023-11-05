package users

import (
	"context"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/tuan882612/apiutils"
)

type repository struct {
	postgres *pgxpool.Pool
	cache    *redis.Client
}

func NewRepository(pg *pgxpool.Pool, cache *redis.Client) *repository {
	return &repository{postgres: pg, cache: cache}
}

func (r *repository) GetUser(ctx context.Context, userID uuid.UUID) (*User, error) {
	user := &User{}

	row := r.postgres.QueryRow(ctx, GetUserQuery, userID)

	if err := user.Scan(row); err != nil {
		if err == pgx.ErrNoRows {
			return nil, apiutils.NewErrNotFound("user not found")
		}

		return nil, err
	}

	return user, nil
}

func (r *repository) VerifyCliKey(ctx context.Context, userID uuid.UUID) (string, error) {
	key := "clikey:" + userID.String()
	data, err := r.cache.Get(key).Result()
	if err != nil {
		if err == redis.Nil {
			return data, apiutils.NewErrNotFound("clikey not found")
		}

		log.Error().Str("location", "VerifyCliKey").Msgf("%v: %v", userID, err)
		return data, err
	}

	return data, nil
}

func (r *repository) CreateCliKey(ctx context.Context, userID uuid.UUID, cliKey string) error {
	key := "clikey:" + userID.String()
	if err := r.cache.Set(key, cliKey, 12*time.Hour).Err(); err != nil {
		log.Error().Str("location", "CreateCliKey").Msgf("%v: failed to add clikey: %v", userID, err)
		return err
	}

	return nil
}
