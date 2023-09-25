package auth

import (
	"encoding/json"
	"io"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/tuan882612/apiutils/securityutils"

	"project/internal/auth/email"
	"project/internal/auth/jwt"
	"project/internal/config"
	"project/internal/database"
)

// Dependencies for the base authentication service.
type Dependencies struct {
	Repository   *Repository // base auth repository
	Service      *Service    // base auth service
	Cache        *Cache      // twofa cache repository
	JWTManager   *jwt.Manager
	EmailManager *email.Manager
}

// Constructor for creating all dependencies for the base authentication service.
func NewDependencies(cfg *config.Configuration) (*Dependencies, error) {
	// initialize data access
	databases, err := database.NewDataAccess(cfg.NumCpu, cfg.PgUrl, cfg.RedisUrl, cfg.RedisPsw)
	if err != nil {
		return nil, err
	}
	repo := NewRepository(databases)
	cache := NewCache(databases)

	// initialize dependencies
	emailManager, err := email.NewManger(cfg)
	if err != nil {
		return nil, err
	}

	service := NewService(repo)
	jwtManager := jwt.NewManager(cfg)

	return &Dependencies{
		Repository:   repo,
		Service:      service,
		Cache:        cache,
		JWTManager:   jwtManager,
		EmailManager: emailManager,
	}, nil
}

// Request data from the login endpoint.
type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,max=32"`
}

func (li *LoginInput) Deserialize(data io.ReadCloser) error {
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

// Request data from the register endpoint.
type RegisterInput struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required,min=16,max=32"`
}

func (ri *RegisterInput) Deserialize(data io.ReadCloser) error {
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

// DTO for RegisterInput to RegisterResp.
type RegisterResp struct {
	UserID uuid.UUID `json:"user_id"`
	RegisterInput
	Registered time.Time `json:"registered"`
	UserStatus string    `json:"user_status"`
}

// Creates a new RegisterResp from the given input and validates it.
func NewRegisterResp(input *RegisterInput) (*RegisterResp, error) {
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
		UserStatus:    "nonreg",
	}
	res.Password = hashedPsw
	return res, nil
}
