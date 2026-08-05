package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"proj/pkg/api/service"
)

func GetIDHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	res, err, status := service.GetTaskID(id)
	if err != nil {
		log.Println(err.Error())
		service.SendErrorJSON(w, "Задача не найдена", status)
		return
	}

	byteResp, err := json.Marshal(res)
	if err != nil {
		service.SendErrorJSON(w, err.Error(), status)
		return
	}

	service.SendOkJSONBytes(w, status, byteResp)

}
