package handler

import (
	"encoding/json"
	"net/http"

	"github.com/lezhnevay/go_final_project_lezhnev/internal/storage"
)

// Обработчик для task/done
func TaskDoneHandler(store storage.Store) http.HandlerFunc { // хэндлер для task/done
	return func(res http.ResponseWriter, req *http.Request) {
		id := req.URL.Query().Get("id")
		err := store.TaskDone(id)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(res).Encode(map[string]string{}); err != nil {
			http.Error(res, `{"error":"Ошибка кодирования JSON"}`, http.StatusInternalServerError)
			return
		}
	}
}
