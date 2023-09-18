package auth

import (
	"encoding/json"
	"errors"
	"io"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (li *LoginInput) Deserialize(data io.ReadCloser) error {
	if err := json.NewDecoder(data).Decode(&li); err != nil {
		log.Error().Str("location", "LoginInput.Deserialize").Msg(err.Error())
		return err
	}

	if li.Email == "" || li.Password == "" {
		msg := "email or password is empty"
		log.Error().Str("location", "LoginInput.Deserialize").Msg(msg)
		return errors.New(msg)
	}

	return nil
}

type RegisterInput struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (ri *RegisterInput) Deserialize(data io.ReadCloser) error {
	if err := json.NewDecoder(data).Decode(&ri); err != nil {
		log.Error().Str("location", "RegisterInput.Deserialize").Msg(err.Error())
		return err
	}

	if ri.Email == "" || ri.Name == "" || ri.Password == "" {
		msg := "email, name, or password is empty"
		log.Error().Str("location", "RegisterInput.Deserialize").Msg(msg)
		return errors.New(msg)
	}

	return nil
}

type RegisterResp struct {
	UserID uuid.UUID `json:"user_id"`
	RegisterInput
	Registered time.Time `json:"registered"`
	UserStatus string    `json:"user_status"`
}

func NewRegisterResp(input *RegisterInput) (*RegisterResp, error) {
	if input == nil {
		msg := "register input is nil"
		log.Error().Str("location", "NewRegisterResp").Msg(msg)
		return nil, errors.New(msg)
	}

	return &RegisterResp{
		UserID:        uuid.New(),
		RegisterInput: *input,
		Registered:    time.Now(),
		UserStatus:    "active",
	}, nil
}

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
