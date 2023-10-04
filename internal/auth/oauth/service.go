package oauth

import (
	"context"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) StartOAuth(ctx context.Context) {

}

func (s *Service) CallbackOAuth(ctx context.Context) {

}
