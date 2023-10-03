package config

import "os"

type ServerConfig struct {
	ApiVersion string `validate:"required"`
	Port       string `validate:"required"`
	GRPCPort   string `validate:"required"`
	Host       string `validate:"required"`
}

func newServerConfig() *ServerConfig {
	return &ServerConfig{
		Host:       os.Getenv("HOST"),
		Port:       os.Getenv("PORT"),
		GRPCPort:   os.Getenv("GRPC_PORT"),
		ApiVersion: os.Getenv("API_VERSION"),
	}
}
