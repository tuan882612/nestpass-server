package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"project/internal/config"
)

// Handles the creation of JWT tokens and other related tasks.
type Manager struct {
	secert   string
	duration time.Duration
}

// Constructor for the JWT manager.
func NewManager(cfg *config.Configuration) *Manager {
	return &Manager{secert: cfg.SignKey, duration: cfg.Duration}
}

// Generates a JWT token for the user.
func (j *Manager) GenerateToken(userID uuid.UUID) (string, error) {
	// create new claims
	claims := NewClaims(userID, j.duration)

	// sign the claims and return the token
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(j.secert))
	if err != nil {
		log.Error().Str("location", "GenerateToken").Msgf("failed to generate token: %v", err)
		return "", err
	}

	return token, nil
}
