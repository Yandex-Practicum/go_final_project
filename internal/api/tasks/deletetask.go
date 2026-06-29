package tasks

import (
	"net/http"

	"final_project/internal/db"
)

func DeleteTaskHandler(res http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	if id == "" {
		writeJSON(res, http.StatusBadRequest, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		writeJSON(res, http.StatusNotFound, map[string]string{"error": "Задача не найдена"})
		return
	}

	writeJSON(res, http.StatusOK, map[string]string{})
}
