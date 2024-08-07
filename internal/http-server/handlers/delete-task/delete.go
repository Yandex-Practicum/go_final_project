package deletetask

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type TaskProvider interface {
	DeleteTask(id string) error
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func New(provider TaskProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		id := r.URL.Query().Get("id")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Отсуетствует id"})
			return
		}
		if id, err := strconv.Atoi(id); err != nil || id < 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "id должен быть положительным числом"})
			return
		}

		err := provider.DeleteTask(id)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(struct{}{})
	}
}
