package handlers

import (
	"log"
	"net/http"
	"proj/pkg/api/service"
)

func PutTaskHandler(w http.ResponseWriter, r *http.Request) {
	byteResp, err, status := service.UpdateTask(w, r)

	if err != nil {
		log.Println(err.Error())
		service.SendErrorJSON(w, err.Error(), status)
		return
	}

	service.SendOkJSONBytes(w, status, byteResp)

}
