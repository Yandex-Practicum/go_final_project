package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	sl "github.com/wissio/go_final_project/internal/lib/logger/slog"
	"github.com/wissio/go_final_project/internal/storage/sqlite"
)

func DeleteTask(logger *slog.Logger, s *sqlite.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.DeleteTask"
		logger = logger.With(slog.String("op", op))

		idQueryParam := r.URL.Query().Get("id")
		if idQueryParam == "" {
			logger.Error("idQueryParam is empty")
			http.Error(w, `{"error":"No task ID specified"}`, http.StatusBadRequest)
			return
		}

		taskID, err := strconv.ParseInt(idQueryParam, 10, 64)
		if err != nil {
			logger.Error("Error parsing task ID", sl.Err(err))
			http.Error(w, `{"error":"Invalid task ID"}`, http.StatusBadRequest)
			return
		}

		err = s.DeleteTask(taskID)
		if errors.Is(err, sqlite.ErrTaskNotFound) {
			logger.Error("Task not found", sl.Err(err))
			http.Error(w, `{"error":"Task not found"}`, http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error("Error deleting task", sl.Err(err))
			http.Error(w, `{"error":"Error deleting task"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(map[string]string{})
		logger.Info("Task deleted successfully")
	}
}
