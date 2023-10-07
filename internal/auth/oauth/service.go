package oauth

import (
	"context"
	"project/internal/config"
)

type Service struct{
	cfg *config.OAuthConfig
}

func NewService(oauthCfg *config.OAuthConfig) *Service {
	return &Service{cfg: oauthCfg}
}

func (s *Service) StartOAuth(ctx context.Context) {

}

func (s *Service) CallbackOAuth(ctx context.Context) {

}
