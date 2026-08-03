package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"proj/pkg/api/service"
)

func GetIDhandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	res, err, status := service.GetTaskID(id)
	if err != nil {
		log.Println(err.Error())
		service.SendErrorJSON(w, "Задача не найдена", status)
		return
	}

	byteJSON, err := json.Marshal(res)
	if err != nil {
		service.SendErrorJSON(w, err.Error(), status)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	w.Write(byteJSON)

}
