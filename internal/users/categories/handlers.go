package categories

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/tuan882612/apiutils"

	"nestpass/internal/dependencies"
	"nestpass/pkg/auth"
	"nestpass/pkg/httputils"
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
	userID, err := auth.UidFromCtx(r.Context())
	if err != nil {
		apiutils.HandleHttpErrors(w, err)
		return
	}

	params := httputils.GetPaginationParams(r)
	categories, err := h.svc.GetAllCategories(r.Context(), userID, params)
	if err != nil {
		apiutils.HandleHttpErrors(w, err)
		return
	}

	resp := apiutils.NewRes(http.StatusOK, "", categories)
	resp.SendRes(w)
}

func (h *Handler) GetCategory(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.UidFromCtx(r.Context())
	if err != nil {
		apiutils.HandleHttpErrors(w, err)
		return
	}

	key := r.URL.Query().Get("key")
	if key == "" {
		apiutils.HandleHttpErrors(w, apiutils.NewErrBadRequest("key is required"))
		return
	}

	category, err := h.svc.GetCategory(r.Context(), userID, key)
	if err != nil {
		apiutils.HandleHttpErrors(w, err)
		return
	}

	resp := apiutils.NewRes(http.StatusOK, "", category)
	resp.SendRes(w)
}

func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	category := &Category{}
	if err := category.Deserialize(r.Body); err != nil {
		apiutils.HandleHttpErrors(w, err)
		return
	}

	categoryResp, err := h.svc.CreateCategory(r.Context(), category)
	if err != nil {
		apiutils.HandleHttpErrors(w, err)
		return
	}

	resp := apiutils.NewRes(http.StatusOK, "", categoryResp)
	resp.SendRes(w)
}

func (h *Handler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	categoryResp := &CategoryResp{}
	if err := categoryResp.Deserialize(r.Body); err != nil {
		apiutils.HandleHttpErrors(w, err)
		return
	}

	category, err := h.svc.UpdateCategory(r.Context(), categoryResp)
	if err != nil {
		apiutils.HandleHttpErrors(w, err)
		return
	}

	resp := apiutils.NewRes(http.StatusOK, "", category)
	resp.SendRes(w)
}

func (h *Handler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.UidFromCtx(r.Context())
	if err != nil {
		apiutils.HandleHttpErrors(w, err)
		return
	}

	categoryIDStr := r.URL.Query().Get("category_id")
	if categoryIDStr == "" {
		apiutils.HandleHttpErrors(w, apiutils.NewErrBadRequest("category_id is required"))
		return
	}

	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		apiutils.HandleHttpErrors(w, apiutils.NewErrBadRequest("invalid category_id"))
		return
	}

	if err := h.svc.DeleteCategory(r.Context(), categoryID, userID); err != nil {
		apiutils.HandleHttpErrors(w, err)
		return
	}

	resp := apiutils.NewRes(http.StatusOK, "", nil)
	resp.SendRes(w)
}
