package oauth

import (
	"net/http"
	"project/internal/config"

	"github.com/tuan882612/apiutils"
)

type Handler struct{
	svc *Service
}

func NewHandler(cfg *config.Configuration) *Handler {
	return &Handler{svc: NewService(cfg.OAuth)}
}

func (h *Handler) Invoke(w http.ResponseWriter, r *http.Request) {
	resp := apiutils.NewRes(http.StatusOK, "", nil)
	resp.SendRes(w)
}

func (h *Handler) Callback(w http.ResponseWriter, r *http.Request) {
	resp := apiutils.NewRes(http.StatusOK, "", nil)
	resp.SendRes(w)
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	resp := apiutils.NewRes(http.StatusOK, "", nil)
	resp.SendRes(w)
}
