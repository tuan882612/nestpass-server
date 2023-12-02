package oauth

import (
	"time"

	"github.com/google/uuid"

	"project/internal/auth"
)

type OAuthData struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (u *OAuthData) NewUser() *auth.Register {
	return &auth.Register{
		UserID:     uuid.New(),
		Email:      u.Email,
		Name:       u.Name,
		Password:   "",
		Registered: time.Now(),
		UserStatus: auth.ActiveUser,
	}
}
