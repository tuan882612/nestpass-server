package singlefa

import (
	"context"

	"project/internal/auth"
)

type Service interface {
	SfaLogin(ctx context.Context, input *auth.LoginInput) (string, error)
	SfaRegister(ctx context.Context, input *auth.RegisterInput) (string, error)
}
