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
func (r *Cache) GetTwofaCache(ctx context.Context, userID uuid.UUID) (*email.TwofaBody, error) {
	data, err := r.cache.Get(userID.String()).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, apiutils.NewErrNotFound("twofa not found")
		}

		log.Error().Str("location", "GetTwofaCache").Msg(err.Error())
		return nil, err
	}

	tfaBody := &email.TwofaBody{}
	if err := tfaBody.Deserialize(data); err != nil {
		log.Error().Str("location", "GetTwofaCache").Msg(err.Error())
		return nil, err
	}

	return tfaBody, nil
}

// Updates the user's twofa data.
func (r *Cache) UpdateTwofaCache(ctx context.Context, userID uuid.UUID, body *email.TwofaBody) error {
	data, err := body.Serialize()
	if err != nil {
		log.Error().Str("location", "UpdateTwofaCache").Msg(err.Error())
		return err
	}

	// update the twofa data and set the ttl to the previous value
	idStr := userID.String()
	ttl := r.cache.TTL(idStr).Val()
	if err := r.cache.Set(idStr, data, ttl).Err(); err != nil {
		log.Error().Str("location", "UpdateTwofaCache").Msg(err.Error())
		return err
	}

	return nil
}

// Deletes the user's twofa data.
func (r *Cache) DeleteTwofaCache(ctx context.Context, userID uuid.UUID) error {
	if err := r.cache.Del(uuid.UUID.String(userID)).Err(); err != nil {
		log.Error().Str("location", "DeleteTwofaCache").Msg(err.Error())
		return err
	}

	return nil
}
