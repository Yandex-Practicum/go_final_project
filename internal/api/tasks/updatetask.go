package tasks

import (
	"encoding/json"
	"net/http"
	"time"

	"final_project/internal/db"
)

func UpdateTaskHandler(res http.ResponseWriter, req *http.Request) {
	var task db.Task

	err := json.NewDecoder(req.Body).Decode(&task)
	if err != nil {
		writeJSON(res, http.StatusBadRequest, map[string]string{"error": "Ошибка десериализации JSON"})
		return
	}
	if task.ID == "" {
		writeJSON(res, http.StatusBadRequest, map[string]string{"error": "Не указан идентификатор"})
		return
	}
	if err := validateAndFormatTask(time.Now(), &task); err != nil {
		writeJSON(res, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	err = db.UpdateTask(&task)
	if err != nil {
		writeJSON(res, http.StatusNotFound, map[string]string{"error": "Задача не найдена"})
		return
	}
	writeJSON(res, http.StatusOK, map[string]string{})
}
