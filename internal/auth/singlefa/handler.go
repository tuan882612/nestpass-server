package singlefa

import (
	"errors"
	"net/http"

	"github.com/rs/zerolog/log"
)

type Handler struct {
	sfaService Service
}

func NewHandler(service Service) (*Handler, error) {
	if service == nil {
		msg := "service is nil"
		log.Error().Str("location", "NewHandler").Msg(msg)
		return nil, errors.New(msg)
	}

	return &Handler{
		sfaService: service,
	}, nil
}

func (s *Handler) Login(w http.ResponseWriter, r *http.Request) {

}

func (s *Handler) Register(w http.ResponseWriter, r *http.Request) {

}
