package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/lezhnevay/go_final_project_lezhnev/configs"
	"github.com/lezhnevay/go_final_project_lezhnev/internal/storage"
)

// Обработчик GET для /tasks
func TasksGetHandler(store storage.Store) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		searchParams := req.URL.Query().Get("search")
		tasks, err := store.GetTasks(searchParams)
		if err != nil {
			err := errors.New("ошибка запроса к базе данных")
			configs.ErrorResponse.Error = err.Error()
			json.NewEncoder(res).Encode(configs.ErrorResponse)
			return
		}
		response := map[string][]configs.Task{
			"tasks": tasks,
		}
		res.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(res).Encode(response); err != nil {
			http.Error(res, `{"error":"Ошибка кодирования JSON"}`, http.StatusInternalServerError)
			return
		}
	}
}
