package helpers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/tuan882612/apiutils"
)

func GetUidHeader(r *http.Request) (uuid.UUID, error) {
	uidStr := r.Header.Get("X-Uid")
	if uidStr == "" {
		return uuid.Nil, apiutils.NewErrBadRequest("missing uid header")
	}

	uid, err := uuid.Parse(uidStr)
	if err != nil {
		return uuid.Nil, apiutils.NewErrBadRequest("invalid uid header")
	}

	return uid, nil
}
