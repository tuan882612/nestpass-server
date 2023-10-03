package database

import (
	"errors"

	"github.com/go-redis/redis"
	"github.com/rs/zerolog/log"
)

func getRedis(redisUrl, redisPsw string) (*redis.Client, error) {
	// check if redis url or password is empty
	if redisUrl == "" || redisPsw == "" {
		errMsg := "redis url or password is empty"
		log.Error().Str("location", "getRedis").Msg(errMsg)
		return nil, errors.New(errMsg)
	}

	log.Info().Msg("connecting redis...")

	// connect to redis database
	conn := redis.NewClient(&redis.Options{
		Addr:     redisUrl,
		Password: redisPsw,
		DB:       0,
	})

	// check if redis is connected
	_, err := conn.Ping().Result()
	if err != nil {
		log.Error().Str("location", "getRedis").Msgf("failed to connect to redis: %v", err)
		return nil, err
	}

	return conn, nil
}
