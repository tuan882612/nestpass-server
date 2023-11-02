package auth

import (
	"context"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/tuan882612/apiutils"

	"project/internal/auth/email"
	"project/internal/database"
)

type CacheType string

const (
	TwoFA      CacheType = "twofa"
	Restricted CacheType = "restricted"
	Session    CacheType = "session"
)

type Cache struct {
	cache *redis.Client
}

func NewCache(databases *database.DataAccess) *Cache {
	return &Cache{cache: databases.Redis}
}

// Generalized Get function for retrieving data from the cache.
func (r *Cache) GetData(ctx context.Context, userID uuid.UUID, mode CacheType) (interface{}, error) {
	if mode == TwoFA {
		// set the keys for the pipeline
		rKey, tfaKey := "restricted:"+userID.String(), "twofa:"+userID.String()

		// create a pipeline
		pipe := r.cache.Pipeline()
		rCmd := pipe.Exists(rKey)
		tfaCmd := pipe.Get(tfaKey)

		// execute the pipeline
		if _, err := pipe.Exec(); err != nil && err != redis.Nil {
			log.Error().Str("location", "GetData.TwoFA").Msgf("%v: failed to get twofa data: %v", userID, err)
			return nil, err
		}

		// check if user is restricted
		if rCmd.Val() == 1 {
			return nil, apiutils.NewErrForbidden("user is restricted")
		}

		// get the user's twofa data
		data, err := tfaCmd.Result()
		if err != nil {
			if err == redis.Nil {
				return nil, apiutils.NewErrNotFound("twofa not found")
			}

			log.Error().Str("location", "GetData.TwoFA").Msgf("%v: failed to get twofa data: %v", userID, err)
			return nil, err
		}

		// deserialize the data into a twofa body
		tfaBody := &email.Twofa{}
		if err := tfaBody.Deserialize(data); err != nil {
			log.Error().Str("location", "GetData.TwoFA").Msgf("%v: failed to deserialize twofa data: %v", userID, err)
			return nil, err
		}

		return tfaBody, nil
	}

	// set the key for either session or restricted
	key := string(mode) + ":" + userID.String()
	_, err := r.cache.Get(key).Result()

	// check for errors
	if err != nil && err != redis.Nil {
		log.Error().Str("location", "GetData.(Restricted, Session)").Msgf("%v: failed to get data: %v", userID, err)
		return nil, err
	}

	// check if session exists or if the user is not restricted
	if err == redis.Nil {
		if mode == Restricted {
			return nil, nil
		}

		return nil, apiutils.NewErrUnauthorized("session expired")
	}

	// checks if the session is active.
	if mode == Session {
		return nil, nil
	}

	return nil, apiutils.NewErrForbidden("user is restricted")
}

// Adds the user as restricted to the cache.
func (r *Cache) AddRestricted(ctx context.Context, userID uuid.UUID) error {
	// set the keys for the pipeline
	rKey, tfaKey := "restricted:"+userID.String(), "twofa:"+userID.String()

	// create a pipeline
	pipe := r.cache.Pipeline()
	pipe.Set(rKey, "restricted", 3*time.Hour)
	pipe.Del(tfaKey)

	// execute the pipeline
	if _, err := pipe.Exec(); err != nil {
		log.Error().Str("location", "AddRestrictedCache").Msgf("%v: failed to add restricted user: %v", userID, err)
		return err
	}

	return nil
}

// Adds a 30 minute ttl session to the cache.
func (r *Cache) AddSession(ctx context.Context, userID uuid.UUID) error {
	key := "session:" + userID.String()
	if err := r.cache.Set(key, "session", 30*time.Minute).Err(); err != nil {
		log.Error().Str("location", "AddSessionCache").Msgf("%v: failed to add session: %v", userID, err)
		return err
	}

	return nil
}

// Adds reset key to the cache.
func (r *Cache) AddResetKey(ctx context.Context, userID uuid.UUID, resetKeyHash string) error {
	key := "reset:" + userID.String()
	if err := r.cache.Set(key, resetKeyHash, 30*time.Minute).Err(); err != nil {
		log.Error().Str("location", "AddResetKeyCache").Msgf("%v: failed to add reset key: %v", userID, err)
		return err
	}
	
	return nil
}

// Updates the user's twofa data.
func (r *Cache) UpdateTwofa(ctx context.Context, userID uuid.UUID, body *email.Twofa) error {
	data, err := body.Serialize()
	if err != nil {
		log.Error().Str("location", "UpdateTwofaCache").Msgf("%v: failed to serialize twofa data: %v", userID, err)
		return err
	}

	// update the twofa data and set the ttl to the previous value
	key := "twofa:" + userID.String()
	ttl := r.cache.TTL(key).Val()
	if err := r.cache.Set(key, data, ttl).Err(); err != nil {
		log.Error().Str("location", "UpdateTwofaCache").Msgf("%v: failed to update twofa data: %v", userID, err)
		return err
	}

	return nil
}

// Deletes user data from the cache.
func (r *Cache) DeleteData(ctx context.Context, userID uuid.UUID, mode CacheType) error {
	key := string(mode) + ":" + userID.String()
	if err := r.cache.Del(key).Err(); err != nil {
		log.Error().Str("location", "DeleteDataCache").Msgf("failed to delete data: %v", err)
		return err
	}

	return nil
}
