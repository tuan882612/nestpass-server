package auth

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"project/internal/auth/email"
)

type Service interface {
	VerifyUser(ctx context.Context, email, password string) (uuid.UUID, error)
	RegisterUser(ctx context.Context, input *RegisterResp) error
}

type Repository interface {
	GetUserCredentials(ctx context.Context, email string) (uuid.UUID, string, error)
	GetTwofaCache(ctx context.Context, userID uuid.UUID) (*email.TwofaBody, error)
	UpdateTwofaCache(ctx context.Context, userID uuid.UUID, body *email.TwofaBody) error
	DeleteTwofaCache(ctx context.Context, userID uuid.UUID) error
	AddUser(ctx context.Context, tx pgx.Tx, input *RegisterResp) error
	UpdateUserStatus(ctx context.Context, userID uuid.UUID) error
	startTx(ctx context.Context) (pgx.Tx, error)
}
