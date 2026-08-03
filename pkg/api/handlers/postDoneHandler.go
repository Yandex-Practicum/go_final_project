package handlers

import (
	"log"
	"net/http"
	"proj/pkg/api/service"
)

func DonePostHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	jsByte, err, status := service.PostDone(id, w, r)
	if err != nil {
		log.Println(err.Error())
		service.SendErrorJSON(w, err.Error(), status)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	w.Write(jsByte)

}
