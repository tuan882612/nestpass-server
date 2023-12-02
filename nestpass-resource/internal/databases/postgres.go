package databases

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func NewPostgres(ctx context.Context, pgUrl string) (*pgxpool.Pool, error) {
	// check if postgres url or numCpu is empty
	if pgUrl == "" {
		errMsg := "postgres url or numCpu is empty"
		log.Error().Str("location", "getPostgres").Msg(errMsg)
		return nil, errors.New(errMsg)
	}

	log.Info().Msg("initializing postgres connection...")

	// parse and set postgres config
	config, err := pgxpool.ParseConfig(pgUrl)
	if err != nil {
		log.Error().Str("location", "getPostgres").Msg(err.Error())
		return nil, err
	}

	// connect to postgres database
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Error().Str("location", "getPostgres").Msg(err.Error())
		return nil, err
	}

	// check if postgres is connected
	if err := pool.Ping(ctx); err != nil {
		log.Error().Str("location", "getPostgres").Msg(err.Error())
		return nil, err
	}

	return pool, nil
}
