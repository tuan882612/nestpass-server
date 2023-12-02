package databases

import (
	"context"

	"github.com/go-redis/redis"
	"github.com/jackc/pgx/v5/pgxpool"

	"nestpass/internal/config"
)

type Databases struct {
	Postgres *pgxpool.Pool
	Redis    *redis.Client
}

func New(cfg *config.Configuration) (*Databases, error) {
	// initialize postgres
	pg, err := NewPostgres(context.Background(), cfg.PgURL)
	if err != nil {
		return nil, err
	}

	// initialize redis
	rds, err := NewRedis(cfg.RedisURL, cfg.RedisPsw)
	if err != nil {
		return nil, err
	}

	return &Databases{
		Postgres: pg,
		Redis:    rds,
	}, nil
}
