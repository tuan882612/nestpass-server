package auth

import (
	"crypto/rand"
	"encoding/json"
	"io"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/tuan882612/apiutils/securityutils"
)

// Base User retrieve model.
type User struct {
	UserID     uuid.UUID
	Password   string
	UserStatus string
}

// user statuses
const (
	NonRegUser   = "nonreg"
	ActiveUser   = "active"
	InactiveUser = "inactive"
)

// Request data from the login endpoint.
type Login struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,max=32"`
}

func (li *Login) Deserialize(data io.ReadCloser) error {
	// deserialize the data
	if err := json.NewDecoder(data).Decode(&li); err != nil {
		log.Error().Str("location", "Login.Deserialize").Msgf("failed to deserialize data: %v", err)
		return err
	}

	// validate the input
	if err := validator.New().Struct(li); err != nil {
		log.Error().Str("location", "Login.Deserialize").Msgf("failed to validate input: %v", err)
		return err
	}

	return nil
}

// Request data from the register endpoint.
type Register struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required,min=16,max=32"`
}

func (ri *Register) Deserialize(data io.ReadCloser) error {
	// deserialize the data
	if err := json.NewDecoder(data).Decode(&ri); err != nil {
		log.Error().Str("location", "Register.Deserialize").Msgf("failed to deserialize data: %v", err)
		return err
	}

	// validate the input
	if err := validator.New().Struct(ri); err != nil {
		log.Error().Str("location", "Register.Deserialize").Msgf("failed to validate input: %v", err)
		return err
	}

	return nil
}

// Request data from the resend code endpoint.
type Resend struct {
	Email string `json:"email" validate:"required,email"`
}

func (ri *Resend) Deserialize(data io.ReadCloser) error {
	// deserialize the data
	if err := json.NewDecoder(data).Decode(&ri); err != nil {
		log.Error().Str("location", "Resend.Deserialize").Msgf("failed to deserialize data: %v", err)
		return err
	}

	// validate the input
	if err := validator.New().Struct(ri); err != nil {
		log.Error().Str("location", "Resend.Deserialize").Msgf("failed to validate input: %v", err)
		return err
	}

	return nil
}

// Request data from the reset password endpoint.
type ResetPsw struct {
	Password string `json:"password" validate:"required,min=16,max=32"`
}

func (rpi *ResetPsw) Deserialize(data io.ReadCloser) error {
	// deserialize the data
	if err := json.NewDecoder(data).Decode(&rpi); err != nil {
		log.Error().Str("location", "ResetPsw.Deserialize").Msgf("failed to deserialize data: %v", err)
		return err
	}

	// validate the input
	if err := validator.New().Struct(rpi); err != nil {
		log.Error().Str("location", "ResetPsw.Deserialize").Msgf("failed to validate input: %v", err)
		return err
	}

	return nil
}

// DTO for Register to RegisterResp.
type RegisterResp struct {
	UserID uuid.UUID `json:"user_id"`
	Register
	Registered time.Time `json:"registered"`
	UserStatus string    `json:"user_status"`
	Salt       []byte    `json:"salt"`
}

// Creates a new RegisterResp from the given input and validates it.
func NewRegisterResp(input *Register) (*RegisterResp, error) {
	// validate the input
	if err := validator.New().Struct(input); err != nil {
		log.Error().Str("location", "NewRegisterResp").Msgf("failed to validate input: %v", err)
		return nil, err
	}

	// hash the password
	hashedPsw, err := securityutils.HashPassword(input.Password)
	if err != nil {
		log.Error().Str("location", "NewRegisterResp").Msgf("failed to hash password: %v", err)
		return nil, err
	}

	// generate a new salt
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}

	// create the new RegisterResp and assign the hashed password
	res := &RegisterResp{
		UserID:     uuid.New(),
		Register:   *input,
		Registered: time.Now(),
		UserStatus: "nonreg",
		Salt:       salt,
	}
	res.Password = hashedPsw
	return res, nil
}
