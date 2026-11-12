package api

import (
	"net/http"
	"time"

	"go_final_project/pkg/nextdate"
)

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.URL.Query().Get("now")
	date := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		http.Error(w, "invalid now date", http.StatusBadRequest)
		return
	}

	result, err := nextdate.NextDate(now, date, repeat)
	if err != nil {
		result = ""
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte(result))
}