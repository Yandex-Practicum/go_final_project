package handler

import (
	"net/http"
	"time"

	"go_final_project/internal/constants"
)

func (h *Handler) NextDate(w http.ResponseWriter, r *http.Request) {
	nowStr := r.URL.Query().Get("now")
	date := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	now, err := time.Parse(constants.DateFormat, nowStr)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	nextDate, err := NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = w.Write([]byte(nextDate))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
