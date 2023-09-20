package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID `json:"user_id"`
}

func NewClaims(userID uuid.UUID, duration time.Duration) (*Claims, error) {
	if userID == uuid.Nil || duration == 0 {
		msg := "invalid claims arguments"
		log.Error().Str("location", "NewClaims").Msg(msg)
		return nil, errors.New(msg)
	}

	return &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "nestpass.auth",
		},
		UserID: userID,
	}, nil
}
