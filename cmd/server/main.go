package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"project/internal/config"
	"project/internal/server"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Logger = log.With().Caller().Logger()

	if err := config.LoadEnv(".env"); err != nil {
		os.Exit(1)
	}

	// initialize and validate new configuration instance
	cfg := config.NewConfiguration()

	if err := cfg.Validate(); err != nil {
		os.Exit(1)
	}

	svr, err := server.New(cfg)
	if err != nil {
		os.Exit(1)
	}

	if err = svr.SetupRouter(); err != nil {
		os.Exit(1)
	}

	svr.Run()
}
