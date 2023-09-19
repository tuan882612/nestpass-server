package auth

import (
	"encoding/json"
	"errors"
	"io"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/tuan882612/apiutils/securityutils"
)

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=16,max=32"`
}

func (li *LoginInput) Deserialize(data io.ReadCloser) error {
	if data == nil {
		msg := "nil data"
		log.Error().Str("location", "LoginInput.Deserialize").Msg(msg)
		return errors.New(msg)
	}

	if err := json.NewDecoder(data).Decode(&li); err != nil {
		log.Error().Str("location", "LoginInput.Deserialize").Msg(err.Error())
		return err
	}

	if err := validator.New().Struct(li); err != nil {
		log.Error().Str("location", "Validate").Msg(err.Error())
		return err
	}

	return nil
}

type RegisterInput struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required,min=16,max=32"`
}

func (ri *RegisterInput) Deserialize(data io.ReadCloser) error {
	if data == nil {
		msg := "nil data"
		log.Error().Str("location", "LoginInput.Deserialize").Msg(msg)
		return errors.New(msg)
	}

	if err := json.NewDecoder(data).Decode(&ri); err != nil {
		log.Error().Str("location", "RegisterInput.Deserialize").Msg(err.Error())
		return err
	}

	if err := validator.New().Struct(ri); err != nil {
		log.Error().Str("location", "Validate").Msg(err.Error())
		return err
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

	if err := validator.New().Struct(input); err != nil {
		log.Error().Str("location", "NewRegisterResp").Msg(err.Error())
		return nil, err
	}

	hashedPsw, err := securityutils.HashPassword(input.Password)
	if err != nil {
		log.Error().Str("location", "NewRegisterResp").Msg(err.Error())
		return nil, err
	}
	
	res := &RegisterResp{
		UserID:        uuid.New(),
		RegisterInput: *input,
		Registered:    time.Now(),
		UserStatus:    "active",
	}
	res.Password = hashedPsw
	return res, nil
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
