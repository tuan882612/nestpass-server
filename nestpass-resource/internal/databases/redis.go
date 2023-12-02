package databases

import (
	"errors"

	"github.com/go-redis/redis"
	"github.com/rs/zerolog/log"
)

func NewRedis(redisUrl, redisPsw string) (*redis.Client, error) {
	// check if redis url or password is empty
	if redisUrl == "" || redisPsw == "" {
		errMsg := "redis url or password is empty"
		log.Error().Str("location", "getRedis").Msg(errMsg)
		return nil, errors.New(errMsg)
	}

	log.Info().Msg("initializing redis connection...")

	// connect to redis database
	conn := redis.NewClient(&redis.Options{
		Addr:     redisUrl,
		Password: redisPsw,
		DB:       0,
	})

	// check if redis is connected
	if err := conn.Ping().Err(); err != nil {
		log.Error().Str("location", "getRedis").Msg(err.Error())
		return nil, err
	}

	return conn, nil
}
