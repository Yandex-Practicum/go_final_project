package handler

import (
	"database/sql"
	"go_final_project/storage"
	"net/http"
	"strconv"
)

// GetTasksHandler возвращает обработчик для получения задач
func GetTasksHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var response []byte
		var responseStatus int
		var err error

		if req.Method != http.MethodGet {
			http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
			return
		}

		id := req.URL.Query().Get("id")
		limitParam := req.URL.Query().Get("limit")
		limit := 50 // Значение по умолчанию
		if limitParam != "" {
			if l, parseErr := strconv.Atoi(limitParam); parseErr == nil && l >= 10 && l <= 50 {
				limit = l
			}
		}

		if id != "" {
			response, responseStatus, err = storage.GetTaskByID(db, id)
		} else {
			response, responseStatus, err = storage.GetTasks(db, limit)
		}

		if err != nil {
			http.Error(w, err.Error(), responseStatus)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}
