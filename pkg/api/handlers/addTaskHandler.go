package handlers

import (
	"log"
	"net/http"
	"proj/pkg/api/service"
)

func AddTaskHandler(w http.ResponseWriter, r *http.Request) {

	byteResp, err, status := service.AddTask(w, r)
	if err != nil {
		log.Println(err.Error())
		service.SendErrorJSON(w, err.Error(), status)
		return
	}

	service.SendOkJSONBytes(w, status, byteResp)

}
