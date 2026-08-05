package handlers

import (
	"net/http"
	"proj/pkg/api/service"
)

func LoginHandle(w http.ResponseWriter, r *http.Request) {
	byteResp, err, status := service.LoginAuther(w, r)
	if err != nil {
		service.SendErrorJSON(w, err.Error(), status)
		return
	}

	service.SendOkJSONBytes(w, status, byteResp)

}
