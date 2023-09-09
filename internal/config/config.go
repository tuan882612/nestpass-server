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
	HOST        string
	PORT        string
	API_VERSION string
	PG_URL      string
	NUM_CPU     int
	REDIS_URL   string
	REDIS_PSW   string
}

func NewConfiguration() *Configuration {
	apiVersion := os.Getenv("API_VERSION")
	port := os.Getenv("PORT")
	host := os.Getenv("HOST")
	pgUrl := os.Getenv("PG_URL")
	redisUrl := os.Getenv("REDIS_URL")
	redisPsw := os.Getenv("REDIS_PSW")

	return &Configuration{
		HOST:        host,
		PORT:        port,
		API_VERSION: apiVersion,
		PG_URL:      pgUrl,
		NUM_CPU:     runtime.NumCPU(),
		REDIS_URL:   redisUrl,
		REDIS_PSW:   redisPsw,
	}
}

func (s *Configuration) Validate() error {
	// create a map of the config values
	configMap := map[string]interface{}{
		"HOST":        s.HOST,
		"PORT":        s.PORT,
		"API_VERSION": s.API_VERSION,
		"PG_URL":      s.PG_URL,
		"NUM_CPU":     s.NUM_CPU,
		"REDIS_URL":   s.REDIS_URL,
		"REDIS_PSW":   s.REDIS_PSW,
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
