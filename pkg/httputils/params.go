package httputils

import (
	"net/http"

	"github.com/google/uuid"
)

type Pagination struct {
	Index string `json:"index"`
	Limit string `json:"limit"`
}

func GetPaginationParams(r *http.Request) *Pagination {
	index := r.URL.Query().Get("index")
	if index == "" {
		index = uuid.Nil.String()
	}

	limit := r.URL.Query().Get("limit")
	if limit == "" {
		limit = "5"
	}

	return &Pagination{
		Index: index,
		Limit: limit,
	}
}
