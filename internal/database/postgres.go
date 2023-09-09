package database

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func getPostgres(ctx context.Context, pgUrl string, numCpu int) (*pgxpool.Pool, error) {
	// check if postgres url or numCpu is empty
	if pgUrl == "" || numCpu == 0 {
		errMsg := "postgres url or numCpu is empty"
		log.Error().Str("location", "getPostgres").Msg(errMsg)
		return nil, errors.New(errMsg)
	}

	log.Info().Msg("connecting postgres...")

	// parse and set postgres config
	config, err := pgxpool.ParseConfig(pgUrl)
	if err != nil {
		log.Error().Str("location", "getPostgres").Msg(err.Error())
		return nil, err
	}

	config.MinConns, config.MaxConns = int32(numCpu)/2, int32(numCpu)*4

	// connect to postgres database
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Error().Str("location", "getPostgres").Msg(err.Error())
		return nil, err
	}

	return pool, nil
}
