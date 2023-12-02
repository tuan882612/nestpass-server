package passwords

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
	repo := NewRepository(deps.Databases.Postgres, deps.Databases.Redis)
	svc := NewService(repo)
	return &Handler{svc: svc}
}

func (h *Handler) RehashAllPasswords(w http.ResponseWriter, r *http.Request) {
	uidStr := r.Header.Get("X-Uid")
	if uidStr == "" {
		apiutils.HandleHttpErrors(w, apiutils.NewErrBadRequest("missing X-Uid header"))
		return
	}

	userID, err := uuid.Parse(uidStr)
	if err != nil {
		apiutils.HandleHttpErrors(w, apiutils.NewErrBadRequest("invalid X-Uid header"))
		return
	}

	if err := h.svc.ReUpdateAllPasswords(r.Context(), userID); err != nil {
		apiutils.HandleHttpErrors(w, err)
		return
	}

	resp := apiutils.NewRes(http.StatusOK, "passwords rehashed", nil)
	resp.SendRes(w)
}

func (h *Handler) GetAllPasswords(w http.ResponseWriter, r *http.Request) {
	pageParams := httputils.GetPaginationParams(r)

	userID, err := auth.UidFromCtx(r.Context())
	if err != nil {
		apiutils.HandleHttpErrors(w, err)
		return
	}

	isGetCategory := true
	categoryID, err := uuid.Parse(r.URL.Query().Get("category_id"))
	if err != nil {
		isGetCategory = false
	}

	var passwords []*Password
	if isGetCategory {
		passwords, err = h.svc.GetAllPasswordsByCategory(r.Context(), userID, categoryID, pageParams)
	} else {
		passwords, err = h.svc.GetAllPasswords(r.Context(), userID, pageParams)
	}

	if err != nil {
		apiutils.HandleHttpErrors(w, err)
		return
	}

	resp := apiutils.NewRes(http.StatusOK, "", passwords)
	resp.SendRes(w)
}

func (h *Handler) GetPassword(w http.ResponseWriter, r *http.Request) {
	pswID, cateryID := r.URL.Query().Get("password_id"), r.URL.Query().Get("category_id")
	if pswID == "" || cateryID == "" {
		apiutils.HandleHttpErrors(w, apiutils.NewErrBadRequest("missing password_id or category_id"))
		return
	}

	passwordID, err := uuid.Parse(pswID)
	if err != nil {
		apiutils.HandleHttpErrors(w, apiutils.NewErrBadRequest("invalid password_id"))
		return
	}

	categoryID, err := uuid.Parse(cateryID)
	if err != nil {
		apiutils.HandleHttpErrors(w, apiutils.NewErrBadRequest("invalid category_id"))
		return
	}

	userID, err := auth.UidFromCtx(r.Context())
	if err != nil {
		apiutils.HandleHttpErrors(w, err)
		return
	}

	password, err := h.svc.GetPassword(r.Context(), passwordID, categoryID, userID)
	if err != nil {
		apiutils.HandleHttpErrors(w, err)
		return
	}

	resp := apiutils.NewRes(http.StatusOK, "", password)
	resp.SendRes(w)
}

func (h *Handler) CreatePassword(w http.ResponseWriter, r *http.Request) {
	psw := &Password{}
	if err := psw.Deserialize(r.Body); err != nil {
		apiutils.HandleHttpErrors(w, err)
		return
	}

	pswID, err := h.svc.CreatePassword(r.Context(), psw)
	if err != nil {
		apiutils.HandleHttpErrors(w, err)
		return
	}

	resp := apiutils.NewRes(http.StatusCreated, "", nil)
	resp.AddHeader(w, map[string]string{"X-Pid": pswID.String()})
	resp.SendRes(w)
}

func (h *Handler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	psw := &Password{}
	if err := psw.Deserialize(r.Body); err != nil {
		apiutils.HandleHttpErrors(w, err)
		return
	}

	if err := h.svc.UpdatePassword(r.Context(), psw); err != nil {
		apiutils.HandleHttpErrors(w, err)
		return
	}

	resp := apiutils.NewRes(http.StatusOK, "", nil)
	resp.SendRes(w)
}

func (h *Handler) DeletePassword(w http.ResponseWriter, r *http.Request) {
	pswID, catID := r.Header.Get("X-Pid"), r.Header.Get("X-Cid")

	passwordID, err := uuid.Parse(pswID)
	if err != nil {
		apiutils.HandleHttpErrors(w, apiutils.NewErrBadRequest("invalid password_id"))
		return
	}

	categoryID, err := uuid.Parse(catID)
	if err != nil {
		apiutils.HandleHttpErrors(w, apiutils.NewErrBadRequest("invalid category_id"))
		return
	}

	userID, err := auth.UidFromCtx(r.Context())
	if err != nil {
		apiutils.HandleHttpErrors(w, err)
		return
	}

	if err := h.svc.DeletePassword(r.Context(), userID, passwordID, categoryID); err != nil {
		apiutils.HandleHttpErrors(w, err)
		return
	}

	resp := apiutils.NewRes(http.StatusOK, "", nil)
	resp.SendRes(w)
}
