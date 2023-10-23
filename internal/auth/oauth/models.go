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

func (u *OAuthData) NewUser() *auth.RegisterResp {
	return &auth.RegisterResp{
		UserID: uuid.New(),
		Register: auth.Register{
			Email:    u.Email,
			Name:     u.Name,
			Password: "",
		},
		Registered: time.Now(),
		UserStatus: auth.ActiveUser,
	}
}
