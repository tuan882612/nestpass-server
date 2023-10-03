package database

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func getPostgres(ctx context.Context, pgUrl string) (*pgxpool.Pool, error) {
	// check if postgres url or numCpu is empty
	if pgUrl == "" {
		errMsg := "postgres url is empty"
		log.Error().Str("location", "getPostgres").Msg(errMsg)
		return nil, errors.New(errMsg)
	}

	log.Info().Msg("connecting postgres...")

	// parse and set postgres config
	config, err := pgxpool.ParseConfig(pgUrl)
	if err != nil {
		log.Error().Str("location", "getPostgres").Msgf("failed to parse postgres config: %v", err)
		return nil, err
	}

	// connect to postgres database
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Error().Str("location", "getPostgres").Msgf("failed to connect to postgres: %v", err)
		return nil, err
	}

	return pool, nil
}
