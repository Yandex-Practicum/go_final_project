package api

import (
	"encoding/json"
	"net/http"

	"github.com/ZhakiyaK/final_project/pkg/db"
)

// Обновление задачи
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {

	var task db.Task

	//Декодирование
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJsonError(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	// Проверка наличия заголовка
	if task.Title == "" {
		writeJsonError(w, "task title required", http.StatusBadRequest)
		return
	}

	// проверка даты
	if err := checkDate(&task); err != nil {
		writeJsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := db.UpdateTask(&task); err != nil {

		writeJsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJson(w, map[string]string{})
}
