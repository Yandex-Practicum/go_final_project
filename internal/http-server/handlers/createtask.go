package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	sl "github.com/wissio/go_final_project/internal/lib/logger/slog"
	"github.com/wissio/go_final_project/internal/models"
	"github.com/wissio/go_final_project/internal/services"
	"github.com/wissio/go_final_project/internal/storage/sqlite"
)

func CreateTask(log *slog.Logger, s *sqlite.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.AddTask"
		log = log.With(slog.String("op", op))

		var task models.Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			log.Error("Error decoding request body", sl.Err(err))
			http.Error(w, `{"error":"JSON deserialization error"}`, http.StatusBadRequest)
			return
		}

		if err := services.ValidateTask(log, &task); err != nil {
			log.Error("Validation error", sl.Err(err))
			http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
			return
		}

		now := time.Now()
		nowStr := now.Format(DateLayout)

		if task.Date == "" {
			task.Date = nowStr
		} else {
			parsedDate, err := time.Parse(DateLayout, task.Date)
			if err != nil || (task.Date != nowStr && parsedDate.Before(now)) {
				if task.Repeat == "" {
					task.Date = nowStr
				} else {
					nextDate, err := services.NextDate(now, task.Date, task.Repeat)
					if err != nil {
						log.Error("Error calculating next date", sl.Err(err))
						http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
						return
					}
					task.Date = nextDate
				}
			}
		}

		id, err := s.CreateTask(&task)
		if err != nil {
			log.Error("Error creating task in database", sl.Err(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp, _ := json.Marshal(map[string]string{"id": strconv.Itoa(id)})
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(resp)
	}
}
