package users

import (
	"net/http"

	"github.com/tuan882612/apiutils"

	"nestpass/internal/dependencies"
	"nestpass/pkg/auth"
)

type Handler struct {
	svc *service
}

func NewHandler(deps *dependencies.Dependencies) *Handler {
	repo := NewRepository(deps.Databases.Postgres, deps.Databases.Redis)
	svc := NewService(repo)
	return &Handler{svc: svc}
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := auth.UidFromCtx(ctx)
	if err != nil {
		apiutils.HandleHttpErrors(w, err)
		return
	}

	user, err := h.svc.GetUser(ctx, userID)
	if err != nil {
		apiutils.HandleHttpErrors(w, err)
		return
	}

	resp := apiutils.NewRes(http.StatusOK, "", user)
	resp.SendRes(w)
}

func (h *Handler) CreateCliKey(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := auth.UidFromCtx(ctx)
	if err != nil {
		apiutils.HandleHttpErrors(w, err)
		return
	}

	key, err := h.svc.CreateCliKey(ctx, userID)
	if err != nil {
		apiutils.HandleHttpErrors(w, err)
		return
	}

	resp := apiutils.NewRes(http.StatusCreated, "", key)
	resp.SendRes(w)
}
