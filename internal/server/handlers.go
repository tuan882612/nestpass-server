package server

import (
	"net/http"

	"github.com/tuan882612/apiutils"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	info := map[string]string{"status": "ok"}
	resp := apiutils.NewRes(http.StatusOK, "", info)
	resp.SendRes(w)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	errMsg := "this resource route does not exist"
	resp := apiutils.NewRes(http.StatusNotFound, errMsg, nil)
	resp.SendRes(w)
}
