package config

import (
	"errors"
	"fmt"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func LoadEnv(path string) error {
	if err := godotenv.Load(path); err != nil {
		errMsg := fmt.Sprintf("failed to load environment variables: %v", err)
		log.Error().Str("location", "LoadEnv").Msg(errMsg)
		return errors.New(errMsg)
	}

	return nil
}
