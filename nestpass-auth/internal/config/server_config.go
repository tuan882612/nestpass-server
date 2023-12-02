package config

import (
	"os"

	"github.com/rs/zerolog/log"
)

type ServerConfig struct {
	ApiVersion string `validate:"required"`
	Port       string `validate:"required"`
	GRPCPort   string `validate:"required"`
	Host       string `validate:"required"`
	ProdEnv    bool
}

func newServerConfig() *ServerConfig {
	rawProdEnv := os.Getenv("PROD_ENV")
	prodEnv := rawProdEnv == "true"

	if !prodEnv {
		if rawProdEnv == "" {
			log.Warn().Msg("PROD_ENV is not set, defaulting to development mode...")
		} else {
			log.Info().Msg("running in development mode...")
		}
	} else {
		log.Info().Msg("running in production mode...")
	}

	return &ServerConfig{
		Host:       os.Getenv("HOST"),
		Port:       os.Getenv("PORT"),
		GRPCPort:   os.Getenv("GRPC_PORT"),
		ApiVersion: os.Getenv("API_VERSION"),
		ProdEnv:    prodEnv,
	}
}
