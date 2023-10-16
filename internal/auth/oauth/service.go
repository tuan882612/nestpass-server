package oauth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"project/internal/config"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Service struct {
	cfg         *config.OAuthConfig
	oauthConfig *oauth2.Config
}

func NewService(oauthCfg *config.OAuthConfig) *Service {
	conf := &oauth2.Config{
		ClientID:     oauthCfg.ClientID,
		ClientSecret: oauthCfg.ClientSecret,
		RedirectURL:  "http://localhost:2001/api/v1/oauth/callback",
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint:     google.Endpoint,
	}
	return &Service{cfg: oauthCfg, oauthConfig: conf}
}

// Generates a CSRF state token
func (s *Service) generateStateOauthCookie() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func (s *Service) getOAuthData(ctx context.Context, token *oauth2.Token) {
	return
}

func (s *Service) StartOAuth(ctx context.Context, state string) string {
	return s.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (s *Service) CallbackOAuth(ctx context.Context, code string, state string, expectedState string) (*oauth2.Token, error) {
	if state != expectedState {
		return nil, fmt.Errorf("invalid oauth state")
	}

	token, err := s.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s *Service) UserLoginSignup(ctx context.Context, token *oauth2.Token) {
	return
}
