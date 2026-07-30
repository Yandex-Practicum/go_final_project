package handlers

import (
	"net/http"
	"proj/pkg/api/service"
)

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := service.GetTask(50, w, r) // в параметре максимальное количество записей
	if err != nil {
		service.SendErrorJSON(w, "Не удалось получить задачи", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(200)
	w.Write(tasks)

}
