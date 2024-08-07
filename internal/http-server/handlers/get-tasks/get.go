package gettasks

import (
	"cactus3d/go_final_project/internal/models"
	"encoding/json"
	"net/http"
	"strconv"
)

type TaskProvider interface {
	GetTasks(search string, offset, limit int) ([]models.Task, error)
}

type Response struct {
	Tasks []models.Task `json:"tasks"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func New(provider TaskProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		var offset, limit int
		var err error

		offs := r.URL.Query().Get("now")
		if offs != "" {
			offset, err = strconv.Atoi(offs)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
				return
			}
		}
		lims := r.URL.Query().Get("date")
		if lims != "" {
			limit, err = strconv.Atoi(lims)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
				return
			}
		}

		search := r.URL.Query().Get("search")

		tasks, err := provider.GetTasks(search, offset, limit)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Response{Tasks: tasks})
	}
}
