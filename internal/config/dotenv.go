package config

import (
	"errors"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func LoadEnv(path string) error {
	if err := godotenv.Load(path); err != nil {
		errMsg := "error loading .env file"
		log.Error().Str("location", "LoadEnv").Msg(errMsg)
		return errors.New(errMsg)
	}

	return nil
}
