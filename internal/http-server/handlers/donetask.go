package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	sl "github.com/wissio/go_final_project/internal/lib/logger/slog"
	"github.com/wissio/go_final_project/internal/services"
	"github.com/wissio/go_final_project/internal/storage/sqlite"
)

func DoneTask(log *slog.Logger, s *sqlite.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.DoneTask"
		log = log.With(slog.String("op", op))

		taskID := r.URL.Query().Get("id")
		if taskID == "" {
			log.Error("No task ID provided")
			http.Error(w, `{"error":"No task ID provided"}`, http.StatusBadRequest)
			return
		}

		id, err := strconv.ParseInt(taskID, 10, 64)
		if err != nil {
			log.Error("Invalid task ID", sl.Err(err))
			http.Error(w, `{"error":"Invalid task ID"}`, http.StatusBadRequest)
			return
		}

		task, err := s.GetTask(taskID)
		if err != nil {
			if errors.Is(err, sqlite.ErrTaskNotFound) {
				log.Error("Task not found", slog.String("id", taskID))
				http.Error(w, `{"error":"Task not found"}`, http.StatusNotFound)
			} else {
				log.Error("Error getting task", sl.Err(err))
				http.Error(w, `{"error":"Error getting task"}`, http.StatusInternalServerError)
			}
			return
		}

		if task.Repeat == "" {
			err = s.DeleteTask(id)
			if err != nil {
				if errors.Is(err, sqlite.ErrTaskNotFound) {
					log.Error("Task not found", slog.String("id", taskID))
					http.Error(w, `{"error":"Task not found"}`, http.StatusNotFound)
				} else {
					log.Error("Error deleting task", sl.Err(err))
					http.Error(w, `{"error":"Error deleting task"}`, http.StatusInternalServerError)
				}
				return
			}
		} else {
			nextDate, err := services.NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				log.Error("Error calculating next date", sl.Err(err))
				http.Error(w, `{"error":"Error calculating next date"}`, http.StatusBadRequest)
				return
			}
			task.Date = nextDate
			err = s.UpdateTask(task)
			if err != nil {
				log.Error("Error updating task", sl.Err(err))
				http.Error(w, `{"error":"Error updating task"}`, http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]string{}); err != nil {
			log.Error("Error encoding response", sl.Err(err))
			http.Error(w, `{"error":"Error encoding response"}`, http.StatusInternalServerError)
		}
	}
}
