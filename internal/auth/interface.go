package auth

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Service interface {
	VerifyUser(ctx context.Context, email, password string) (uuid.UUID, error)
	RegisterUser(ctx context.Context, input *RegisterInput) (uuid.UUID, error)
}

type Repository interface {
	GetUserCredentials(ctx context.Context, email string) (uuid.UUID, string, error)
	AddUser(ctx context.Context, tx pgx.Tx, input *RegisterResp) error
	startTx(ctx context.Context) (pgx.Tx, error)
}
