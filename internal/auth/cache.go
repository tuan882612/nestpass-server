package auth

import (
	"context"

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
	data, err := r.cache.Get(userID.String()).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, apiutils.NewErrNotFound("twofa not found")
		}

		log.Error().Str("location", "GetTwofaCache").Msgf("%v: failed to get twofa data: %v", userID, err)
		return nil, err
	}

	tfaBody := &email.TwofaBody{}
	if err := tfaBody.Deserialize(data); err != nil {
		log.Error().Str("location", "GetTwofaCache").Msgf("%v: failed to deserialize twofa data: %v", userID, err)
		return nil, err
	}

	return tfaBody, nil
}

// Updates the user's twofa data.
func (r *Cache) UpdateTwofa(ctx context.Context, userID uuid.UUID, body *email.TwofaBody) error {
	data, err := body.Serialize()
	if err != nil {
		log.Error().Str("location", "UpdateTwofaCache").Msgf("%v: failed to serialize twofa data: %v", userID, err)
		return err
	}

	// update the twofa data and set the ttl to the previous value
	idStr := userID.String()
	ttl := r.cache.TTL(idStr).Val()
	if err := r.cache.Set(idStr, data, ttl).Err(); err != nil {
		log.Error().Str("location", "UpdateTwofaCache").Msgf("%v: failed to update twofa data: %v", userID, err)
		return err
	}

	return nil
}

// Deletes the user's twofa data.
func (r *Cache) DeleteTwofa(ctx context.Context, userID uuid.UUID) error {
	if err := r.cache.Del(uuid.UUID.String(userID)).Err(); err != nil {
		log.Error().Str("location", "DeleteTwofaCache").Msgf("%v: failed to delete twofa data: %v", userID, err)
		return err
	}

	return nil
}
