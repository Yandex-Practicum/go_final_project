package handlers

import (
	"log"
	"net/http"
	"proj/pkg/api/service"
)

func PutTaskHandler(w http.ResponseWriter, r *http.Request) {
	ByteResp, err, status := service.UpdateTask(w, r)

	if err != nil {
		log.Println(err.Error())
		service.SendErrorJSON(w, err.Error(), status)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	w.Write(ByteResp)
}
