package api

import (
	"net/http"
	"time"

	"go_final_project/internal/config"
	"go_final_project/internal/scheduler"
)

// NextDateHandler обрабатывает запросы к /api/nextdate
func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.URL.Query().Get("now")
	dateStr := r.URL.Query().Get("date")
	repeatStr := r.URL.Query().Get("repeat")

	now, err := time.Parse(config.TimeFormat, nowStr)
	if err != nil {
		http.Error(w, `{"error": "недопустимая дата now"}`, http.StatusBadRequest)
		return
	}

	nextDate, err := scheduler.NextDate(now, dateStr, repeatStr)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(nextDate))
}
