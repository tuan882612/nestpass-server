package config

import (
	"os"
	"time"
)

type JWTConfig struct {
	Duration time.Duration `validate:"required"`
	SignKey  string        `validate:"required"`
	EmailKey string        `validate:"required"`
}

func newJWTConfig() *JWTConfig {
	duration, err := time.ParseDuration(os.Getenv("TOKEN_DURATION"))
	if err != nil {
		duration = 12 * time.Hour
	}

	return &JWTConfig{
		SignKey:  os.Getenv("SIGN_KEY"),
		EmailKey: os.Getenv("EMAIL_KEY"),
		Duration: duration,
	}
}
