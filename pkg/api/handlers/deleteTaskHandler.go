package handlers

import (
	"log"
	"net/http"
	"proj/pkg/api/service"
)

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	byteResp, err, status := service.DeleteTask(id)
	if err != nil {
		log.Println(err.Error())
		service.SendErrorJSON(w, err.Error(), status)
		return
	}

	service.SendOkJSONBytes(w, status, byteResp)

}
