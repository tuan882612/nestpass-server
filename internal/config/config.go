package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

type Configuration struct {
	Server   *ServerConfig
	Database *DatabaseConfig
	JWT      *JWTConfig
	OAuth    *OAuthConfig
}

func New() *Configuration {
	return &Configuration{
		Server:   newServerConfig(),
		Database: newDatabaseConfig(),
		JWT:      newJWTConfig(),
		OAuth:    newOauthConfig(),
	}
}

func (c *Configuration) Validate() error {
	configMap := map[string]interface{}{
		"Server":   c.Server,
		"Database": c.Database,
		"JWT":      c.JWT,
		"OAuth":    c.OAuth,
	}

	validator := validator.New()
	for name, config := range configMap {
		if err := validator.Struct(config); err != nil {
			log.Error().Str("location", "Configuration.Validate").Msgf("failed to validate %s config: %v", name, err)
			return err
		}
	}

	return nil
}
