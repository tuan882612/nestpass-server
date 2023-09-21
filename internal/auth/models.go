package auth

import (
	"encoding/json"
	"errors"
	"io"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/tuan882612/apiutils/securityutils"

	"project/internal/auth/jwt"
	"project/internal/config"
)

// This struct that holds the dependencies for the base authentication service.
type Deps struct {
	Service    Service // base authentication service
	JWTManger *jwt.JWTManger
}

func NewDependencies(cfg *config.Configuration) (*Deps, error) {
	if cfg == nil {
		msg := "nil configuration"
		log.Error().Str("location", "NewDependencies").Msg(msg)
		return nil, errors.New(msg)
	}

	// initialize sub-dependency repository
	repo, err := NewRepository(cfg)
	if err != nil {
		return nil, err
	}

	// initialize base authentication service
	service, err := NewService(repo)
	if err != nil {
		return nil, err
	}

	// initialize JWT manager
	jwtManager, err := jwt.NewJWTManager(cfg.SignKey, cfg.Duration)
	if err != nil {
		return nil, err
	}

	return &Deps{
		Service:    service,
		JWTManger: jwtManager,
	}, nil
}

// This struct is used to take in request data from the login endpoint.
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

	// deserialize the data
	if err := json.NewDecoder(data).Decode(&li); err != nil {
		log.Error().Str("location", "LoginInput.Deserialize").Msg(err.Error())
		return err
	}

	// validate the input
	if err := validator.New().Struct(li); err != nil {
		log.Error().Str("location", "Validate").Msg(err.Error())
		return err
	}

	return nil
}

// This struct is used to take in request data from the register endpoint.
type RegisterInput struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required,max=32"`
}

func (ri *RegisterInput) Deserialize(data io.ReadCloser) error {
	if data == nil {
		msg := "nil data"
		log.Error().Str("location", "LoginInput.Deserialize").Msg(msg)
		return errors.New(msg)
	}

	// deserialize the data
	if err := json.NewDecoder(data).Decode(&ri); err != nil {
		log.Error().Str("location", "RegisterInput.Deserialize").Msg(err.Error())
		return err
	}

	// validate the input
	if err := validator.New().Struct(ri); err != nil {
		log.Error().Str("location", "Validate").Msg(err.Error())
		return err
	}

	return nil
}

// This struct is used to convert RegisterInput to RegisterResp and validate the input.
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

	// validate the input
	if err := validator.New().Struct(input); err != nil {
		log.Error().Str("location", "NewRegisterResp").Msg(err.Error())
		return nil, err
	}

	// hash the password
	hashedPsw, err := securityutils.HashPassword(input.Password)
	if err != nil {
		log.Error().Str("location", "NewRegisterResp").Msg(err.Error())
		return nil, err
	}

	// create the new RegisterResp and assign the hashed password
	res := &RegisterResp{
		UserID:        uuid.New(),
		RegisterInput: *input,
		Registered:    time.Now(),
		UserStatus:    "active",
	}
	res.Password = hashedPsw
	return res, nil
}
