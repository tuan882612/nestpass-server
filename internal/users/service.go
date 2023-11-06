package users

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"
)

type service struct {
	repo *repository
}

func NewService(repo *repository) *service {
	return &service{repo: repo}
}

func (s *service) GetUser(ctx context.Context, userID uuid.UUID) (*User, error) {
	return s.repo.GetUser(ctx, userID)
}

func (s *service) CreateCliKey(ctx context.Context, userID uuid.UUID) (string, error) {
	cliKeyBytes := make([]byte, 32)
	if _, err := rand.Read(cliKeyBytes); err != nil {
		return "", fmt.Errorf("failed to generate CLI key: %w", err)
	}
	cliKey := hex.EncodeToString(cliKeyBytes)

	if err := s.repo.CreateCliKey(ctx, userID, cliKey); err != nil {
		return "", err
	}

	return cliKey, nil
}
