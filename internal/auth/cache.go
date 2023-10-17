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

type Cache struct {
	cache *redis.Client
}

func NewCache(databases *database.DataAccess) *Cache {
	return &Cache{cache: databases.Redis}
}

// Retrieves the user's twofa data.
func (r *Cache) GetTwofa(ctx context.Context, userID uuid.UUID) (*email.TwofaBody, error) {
	key := "twofa:" + userID.String()
	data, err := r.cache.Get(key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, apiutils.NewErrNotFound("twofa not found")
		}

		log.Error().Str("location", "GetTwofaCache").Msgf("%v: failed to get twofa data: %v", userID, err)
		return nil, err
	}

	// check if data is restricted
	if err := r.GetRestricted(ctx, userID); err != nil {
		return nil, apiutils.NewErrForbidden("user is restricted")
	}

	tfaBody := &email.TwofaBody{}
	if err := tfaBody.Deserialize(data); err != nil {
		log.Error().Str("location", "GetTwofaCache").Msgf("%v: failed to deserialize twofa data: %v", userID, err)
		return nil, err
	}

	return tfaBody, nil
}

// Checks if the user is restricted.
func (r *Cache) GetRestricted(ctx context.Context, userID uuid.UUID) error {
	key := "restricted:" + userID.String()
	_, err := r.cache.Get(key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		
		log.Error().Str("location", "IsRestrictedCache").Msgf("%v: failed to get restricted user: %v", userID, err)
		return err
	}

	return apiutils.NewErrForbidden("user is restricted")
}

// Get the user's session.
func (r *Cache) GetSession(ctx context.Context, userID uuid.UUID) error {
	key := "session:" + userID.String()
	_, err := r.cache.Get(key).Result()
	if err != nil {
		if err == redis.Nil {
			return apiutils.NewErrUnauthorized("session expired")
		}

		log.Error().Str("location", "GetSessionCache").Msgf("%v: failed to get session: %v", userID, err)
		return err
	}

	return nil
}

// Adds the user as restricted to the cache.
func (r *Cache) AddRestricted(ctx context.Context, userID uuid.UUID) error {
	key := "restricted:" + userID.String()
	if err := r.cache.Set(key, "restricted", 3*time.Hour).Err(); err != nil {
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

// Updates the user's twofa data.
func (r *Cache) UpdateTwofa(ctx context.Context, userID uuid.UUID, body *email.TwofaBody) error {
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

// Deletes the user's twofa data.
func (r *Cache) DeleteTwofa(ctx context.Context, userID uuid.UUID) error {
	key := "twofa:" + userID.String()
	if err := r.cache.Del(key).Err(); err != nil {
		log.Error().Str("location", "DeleteTwofaCache").Msgf("%v: failed to delete twofa data: %v", userID, err)
		return err
	}

	return nil
}
