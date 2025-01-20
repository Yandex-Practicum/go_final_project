package handler

import (
	"encoding/json"
	"net/http"

	"github.com/lezhnevay/go_final_project_lezhnev/configs"
	"github.com/lezhnevay/go_final_project_lezhnev/internal/storage"
)

// Обработчик POST для /task
func TaskPostHandler(store storage.Store) http.HandlerFunc { // хэндлер для создания задачи
	return func(res http.ResponseWriter, req *http.Request) {
		var t configs.Task
		err := json.NewDecoder(req.Body).Decode(&t)
		if err != nil {
			http.Error(res, `{"error":"Ошибка десериализации JSON"}`, http.StatusBadRequest)
			return
		}
		id, err := store.CreateTask(t)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		response := configs.Response{ID: id}

		res.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(res).Encode(response); err != nil {
			http.Error(res, `{"error":"Ошибка кодирования JSON"}`, http.StatusInternalServerError)
			return
		}
	}
}
