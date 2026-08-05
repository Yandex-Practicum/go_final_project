package handlers

import (
	"log"
	"net/http"
	"proj/pkg/api/service"
)

func DonePostHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	byteResp, err, status := service.PostDone(id, w, r)
	if err != nil {
		log.Println(err.Error())
		service.SendErrorJSON(w, err.Error(), status)
	}
	service.SendOkJSONBytes(w, status, byteResp)

}
