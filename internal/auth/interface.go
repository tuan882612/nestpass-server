package auth

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Service interface {
	VerifyUser(ctx context.Context, email, password string) (string, error)
	RegisterUser(ctx context.Context, input RegisterInput) (string, error)
}

type Repository interface {
	GetUserCredentials(ctx context.Context, email string) (uuid.UUID, string, error)
	AddUser(ctx context.Context, tx pgx.Tx, input *RegisterResp) error
}
