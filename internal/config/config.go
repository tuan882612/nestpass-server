package config

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/rs/zerolog/log"
)

type Configuration struct {
	Host       string
	Port       string
	ApiVersion string
	PgUrl      string
	NumCpu     int
	RedisUrl   string
	RedisPsw   string
	SignKey    string
	Duration   time.Duration
}

func NewConfiguration() *Configuration {
	apiVersion := os.Getenv("API_VERSION")
	port := os.Getenv("PORT")
	host := os.Getenv("HOST")
	pgUrl := os.Getenv("PG_URL")
	redisUrl := os.Getenv("REDIS_URL")
	redisPsw := os.Getenv("REDIS_PSW")
	singKey := os.Getenv("SIGN_KEY")
	duration, err := time.ParseDuration(os.Getenv("TOKEN_DURATION"))
	if err != nil {
		duration = 12
	}

	return &Configuration{
		Host:       host,
		Port:       port,
		ApiVersion: apiVersion,
		PgUrl:      pgUrl,
		NumCpu:     runtime.NumCPU(),
		RedisUrl:   redisUrl,
		RedisPsw:   redisPsw,
		SignKey:    singKey,
		Duration:   duration,
	}
}

func (s *Configuration) Validate() error {
	// create a map of the config values
	configMap := map[string]interface{}{
		"HOST":       s.Host,
		"PORT":       s.Port,
		"ApiVersion": s.ApiVersion,
		"PG_URL":     s.PgUrl,
		"NUM_CPU":    s.NumCpu,
		"REDIS_URL":  s.RedisUrl,
		"REDIS_PSW":  s.RedisPsw,
		"SIGN_KEY":   s.SignKey,
		"DURATION":   s.Duration,
	}

	// check if any of the config values are empty
	emptyKeys := []string{}
	for key, value := range configMap {
		switch v := value.(type) {
		case string:
			if v == "" {
				emptyKeys = append(emptyKeys, key)
			}
		case int, time.Duration:
			if v == 0 {
				emptyKeys = append(emptyKeys, key)
			}
		}
	}

	// return an error if any of the config values are empty
	if len(emptyKeys) > 0 {
		errMsg := fmt.Sprintf("%v are required", emptyKeys)
		log.Error().Str("location", "Validate").Msg(errMsg)
		return errors.New(errMsg)
	}

	return nil
}
