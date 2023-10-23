package categories

import (
	"net/http"

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

func (h *Handler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	resp := apiutils.NewRes(http.StatusOK, "success", nil)
	resp.SendRes(w)
}

func (h *Handler) GetCategory(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) UpdateCategory(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) DeleteCategory(w http.ResponseWriter, r *http.Request) {

}
