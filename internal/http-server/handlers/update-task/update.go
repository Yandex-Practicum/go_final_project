package updatetask

import (
	"cactus3d/go_final_project/internal/models"
	"cactus3d/go_final_project/internal/utils"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type TaskProvider interface {
	UpdateTask(*models.Task) error
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func New(provider TaskProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req models.Task
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Некорректный формат запроса"})
			return
		}

		if req.Id == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Отсуетствует id"})
			return
		}
		if id, err := strconv.Atoi(req.Id); err != nil || id < 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "id должен быть положительным числом"})
			return
		}

		if req.Title == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Не указан заголовок задачи"})
			return
		}

		if req.Date == "" {
			req.Date = time.Now().Format("20060102")
		} else {
			_, err = time.Parse("20060102", req.Date)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(ErrorResponse{Error: "Неверный формат времени"})
				return
			}
		}

		if req.Repeat != "" {
			_, err = utils.NextDate(time.Now(), req.Date, req.Repeat)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(ErrorResponse{Error: "Неверный формат повторений"})
				return
			}
		}

		err = provider.UpdateTask(&req)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(struct{}{})
	}
}
