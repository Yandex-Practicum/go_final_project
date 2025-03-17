package handlers

import (
	"log/slog"
	"net/http"
	"time"

	sl "github.com/wissio/go_final_project/internal/lib/logger/slog"
	"github.com/wissio/go_final_project/internal/services"
	"github.com/wissio/go_final_project/internal/storage/sqlite"
)

func NextDate(log *slog.Logger, s *sqlite.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nowStr := r.URL.Query().Get("now")
		dateStr := r.URL.Query().Get("date")
		repeatStr := r.URL.Query().Get("repeat")

		now, err := time.Parse(services.DateLayout, nowStr)
		if err != nil {
			log.Error("Failed to parse time", sl.Err(err), slog.String("time", nowStr))
			http.Error(w, "Invalid time format", http.StatusBadRequest)
			return
		}

		nextDate, err := services.NextDate(now, dateStr, repeatStr)
		if err != nil {
			log.Error("Failed to calculate next date", sl.Err(err),
				slog.String("now", nowStr), slog.String("date", dateStr), slog.String("repeat", repeatStr))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Write([]byte(nextDate))
	}
}
