package twofa

import (
	"context"

	"github.com/google/uuid"

	"project/internal/auth"
	"project/internal/auth/email"
)

type Service interface {
	SendVerificationEmail(ctx context.Context, userID uuid.UUID, email string) error
	LoginSend(ctx context.Context, input *auth.LoginInput) (string, error)
	VerifyAuthToken(ctx context.Context, userID uuid.UUID, input *email.TokenInput) (string, error)
	RegisterSend(ctx context.Context, input *auth.RegisterInput) (string, error)
	RegisterVerify(ctx context.Context, userID uuid.UUID, input *email.TokenInput) (string, error)
}
