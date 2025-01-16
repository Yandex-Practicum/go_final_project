package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/lezhnevay/go_final_project_lezhnev/configs"
	"github.com/lezhnevay/go_final_project_lezhnev/internal/storage"
)

// Обработчик GET для /task
func TaskGetHandler(store storage.Store) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		//var t configs.Task
		id := req.URL.Query().Get("id")
		task, err := store.GetTask(id)
		if err != nil {
			err := errors.New("Задача с таким id не найдена")
			configs.ErrorResponse.Error = err.Error()
			json.NewEncoder(res).Encode(configs.ErrorResponse)
			return
		}
		// Возвращаем ответ
		res.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(res).Encode(task); err != nil {
			http.Error(res, `{"error":"Ошибка кодирования JSON"}`, http.StatusInternalServerError)
			return
		}
	}
}
