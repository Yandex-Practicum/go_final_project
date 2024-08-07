package addtask

import (
	"cactus3d/go_final_project/internal/utils"
	"encoding/json"
	"net/http"
	"time"
)

type TasksProvider interface {
	AddTask(date, title, comment, repeat string) (int, error)
}

type Request struct {
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type Response struct {
	Id int `json:"id"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func New(provider TasksProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req Request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Некорректный формат запроса"})
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

		id, err := provider.AddTask(req.Date, req.Title, req.Comment, req.Repeat)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Response{Id: id})
	}
}
