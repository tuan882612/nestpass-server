package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/tuan882612/apiutils"
)

type ctxKey string

const CtxUserID ctxKey = "user_id"

// Claims represents the JWT claims.
type claims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID `json:"user_id"`
}

// DecodeToken decodes a JWT token and returns the Claims.
func DecodeToken(token, signKey string) (*claims, error) {
	// decode token with the given sign key
	payload, err := jwt.ParseWithClaims(token, &claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(signKey), nil
	})

	// handle all possible errors from parsing the token
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, apiutils.NewErrBadRequest(err.Error())
		} else if errors.Is(err, jwt.ErrSignatureInvalid) {
			return nil, apiutils.NewErrUnauthorized(err.Error())
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, apiutils.NewErrUnauthorized(err.Error())
		}

		log.Error().Str("location", "GetPayload").Msg(err.Error())
		return nil, err
	}

	// check if the decoded token is valid claims
	claims, ok := payload.Claims.(*claims)
	if !ok {
		errMsg := "error parsing claims"
		log.Error().Str("location", "New").Msg(errMsg)
		return nil, errors.New(errMsg)
	}

	return claims, nil
}

// Tries to get JWT token from the authorization header.
func GetBearerToken(r *http.Request) (string, error) {
	// get the token from the authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", apiutils.NewErrUnauthorized("missing authorization header")
	}

	// tries to parse the token from the authorization header
	splitToken := strings.Split(authHeader, "Bearer ")
	if len(splitToken) != 2 || splitToken[1] == "" {
		return "", apiutils.NewErrUnauthorized("invalid authorization header")
	}

	return splitToken[1], nil
}

// Parses and retrieves the user id from the context.
func UidFromCtx(ctx context.Context) (uuid.UUID, error) {
	uid, ok := ctx.Value(CtxUserID).(uuid.UUID)
	if !ok {
		errMsg := "error parsing user id"
		log.Error().Str("location", "UidFromCtx").Msg(errMsg)
		return uuid.Nil, errors.New(errMsg)
	}

	return uid, nil
}
