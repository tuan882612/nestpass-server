package passwords

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tuan882612/apiutils"

	"nestpass/internal/dependencies"
)

type Handler struct {
	svc *service
}

func NewHandler(deps *dependencies.Dependencies) *Handler {
	repo := NewRepository(deps.Databases.Postgres)
	svc := NewService(repo)
	return &Handler{svc: svc}
}

func (h *Handler) GetAllPasswords(w http.ResponseWriter, r *http.Request) {
	cid := chi.URLParam(r, "cid")

	resp := apiutils.NewRes(http.StatusOK, "success", cid)
	resp.SendRes(w)
}

func (h *Handler) GetPassword(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) CreatePassword(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) UpdatePassword(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) DeletePassword(w http.ResponseWriter, r *http.Request) {

}
