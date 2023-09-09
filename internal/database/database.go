package database

import (
	"context"

	"github.com/go-redis/redis"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type DataAccess struct {
	Postgres *pgxpool.Pool
	Redis    *redis.Client
}

func NewDataAccess(nCpu int, pgUrl, rdUrl, rdPsw string) (*DataAccess, error) {
	postgres, err := getPostgres(context.Background(), pgUrl, nCpu)
	if err != nil {
		return nil, err
	}

	redis, err := getRedis(rdUrl, rdPsw)
	if err != nil {
		return nil, err
	}

	log.Info().Msg("initializing data access...")
	return &DataAccess{
		Postgres: postgres,
		Redis:    redis,
	}, nil
}
