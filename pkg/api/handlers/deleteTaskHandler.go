package handlers

import (
	"log"
	"net/http"
	"proj/pkg/api/service"
)

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	byteJSON, err, status := service.DeleteTask(id)
	if err != nil {
		log.Println(err.Error())
		service.SendErrorJSON(w, err.Error(), status)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	w.Write(byteJSON)

}
