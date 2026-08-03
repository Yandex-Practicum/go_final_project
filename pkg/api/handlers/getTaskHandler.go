package handlers

import (
	"log"
	"net/http"
	"proj/pkg/api/service"
)

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err, status := service.GetListTask(50, w, r) // в параметре максимальное количество записей
	if err != nil {
		log.Println(err.Error())
		service.SendErrorJSON(w, "Не удалось получить задачи", status)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	w.Write(tasks)

}
