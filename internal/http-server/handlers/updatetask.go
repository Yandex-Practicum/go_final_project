package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	sl "github.com/wissio/go_final_project/internal/lib/logger/slog"
	"github.com/wissio/go_final_project/internal/models"
	"github.com/wissio/go_final_project/internal/services"
	"github.com/wissio/go_final_project/internal/storage/sqlite"
)

func UpdateTask(log *slog.Logger, s *sqlite.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.UpdateTask"
		log = log.With(slog.String("op", op))

		var task models.Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			log.Error("Error decoding request body", sl.Err(err))
			http.Error(w, `{"error":"Failed to decode request"}`, http.StatusBadRequest)
			return
		}

		if err := services.ValidateTask(log, &task); err != nil {
			log.Error("Validation error", sl.Err(err))
			http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
			return
		}

		if err := s.UpdateTask(task); err != nil {
			log.Error("Error updating task in database", sl.Err(err))
			http.Error(w, `{"error":"Failed to update task in database"}`, http.StatusInternalServerError)
			return
		}

		updatedTask, err := s.GetTask(task.Id)
		if err != nil {
			log.Error("Error fetching updated task from database", sl.Err(err))
			http.Error(w, `{"error":"Failed to fetch updated task from database"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(updatedTask)
	}
}
