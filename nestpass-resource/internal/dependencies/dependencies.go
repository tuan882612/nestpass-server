package dependencies

import (
	"github.com/rs/zerolog/log"

	"nestpass/internal/config"
	"nestpass/internal/databases"
)

// Dependencies contains all dependencies for the server.
type Dependencies struct {
	Databases *databases.Databases
}

// New creates a new dependencies instance.
func New(cfg *config.Configuration) (*Dependencies, error) {
	log.Info().Msg("initializing server dependencies...")

	db, err := databases.New(cfg)
	if err != nil {
		return nil, err
	}

	return &Dependencies{
		Databases: db,
	}, nil
}
