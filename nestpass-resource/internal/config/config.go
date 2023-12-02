package config

import (
	"errors"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

type Configuration struct {
	ApiVersion string `validate:"required"`
	Port       string `validate:"required"`
	Host       string `validate:"required"`
	PgURL      string `validate:"required"`
	RedisURL   string `validate:"required"`
	RedisPsw   string `validate:"required"`
	SignKey    string `validate:"required"`
}

func New() *Configuration {
	apiVersion := os.Getenv("API_VERSION")
	port := os.Getenv("PORT")
	host := os.Getenv("HOST")
	pgUrl := os.Getenv("PG_URL")
	redisUrl := os.Getenv("REDIS_URL")
	redisPsw := os.Getenv("REDIS_PSW")
	singKey := os.Getenv("SIGN_KEY")

	return &Configuration{
		Host:       host,
		Port:       port,
		ApiVersion: apiVersion,
		PgURL:      pgUrl,
		RedisURL:   redisUrl,
		RedisPsw:   redisPsw,
		SignKey:    singKey,
	}
}

func (s *Configuration) Validate() error {
	if err := validator.New().Struct(s); err != nil {
		errMsg := err.Error()
		log.Error().Str("location", "config/validate").Msg(errMsg)
		return errors.New(errMsg)
	}

	return nil
}
