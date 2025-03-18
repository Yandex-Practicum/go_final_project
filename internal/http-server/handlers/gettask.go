package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	sl "github.com/wissio/go_final_project/internal/lib/logger/slog"
	"github.com/wissio/go_final_project/internal/storage/sqlite"
)

func GetTask(log *slog.Logger, s *sqlite.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.GetTask"
		logger := log.With(slog.String("op", op))

		taskID := r.URL.Query().Get("id")
		if taskID == "" {
			http.Error(w, `{"error":"No task ID specified"}`, http.StatusBadRequest)
			return
		}

		task, err := s.GetTask(taskID)
		if err != nil {
			logger.Error("Failed to get task", sl.Err(err))
			http.Error(w, `{"error":"Invalid task ID"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(task); err != nil {
			logger.Error("Failed to write response", sl.Err(err))
		}
	}
}
