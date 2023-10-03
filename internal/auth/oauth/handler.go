package oauth

import (
	"net/http"

	"github.com/tuan882612/apiutils"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	resp := apiutils.NewRes(http.StatusOK, "", nil)
	resp.SendRes(w)
}
