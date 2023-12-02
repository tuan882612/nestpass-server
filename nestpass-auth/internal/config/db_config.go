package config

import "os"

type DatabaseConfig struct {
	PgURL    string `validate:"required"`
	RedisURL string `validate:"required"`
	RedisPsw string `validate:"required"`
}

func newDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		PgURL:    os.Getenv("PG_URL"),
		RedisURL: os.Getenv("REDIS_URL"),
		RedisPsw: os.Getenv("REDIS_PSW"),
	}
}
