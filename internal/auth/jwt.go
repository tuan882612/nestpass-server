package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type JWTManger struct {
	secert   string
	duration time.Duration
}

func NewJWTManager(cfgSecert string, tokenDuration time.Duration) (*JWTManger, error) {
	if cfgSecert == "" || tokenDuration == 0 {
		msg := "invalid jwt handler arguments"
		log.Error().Str("location", "NewJWTManager").Msg(msg)
		return nil, errors.New(msg)
	}

	return &JWTManger{
		secert:   cfgSecert,
		duration: tokenDuration,
	}, nil
}

func (j *JWTManger) GenerateToken(userID uuid.UUID) (string, error) {
	if userID == uuid.Nil {
		msg := "nil userID"
		log.Error().Str("location", "GenerateToken").Msg(msg)
		return "", errors.New(msg)
	}

	claims, err := NewClaims(userID, j.duration)
	if err != nil {
		log.Error().Str("location", "GenerateToken").Msg(err.Error())
		return "", err
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(j.secert))
	if err != nil {
		log.Error().Str("location", "GenerateToken").Msg(err.Error())
		return "", err
	}

	return token, nil
}
