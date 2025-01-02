package handlers

import (
	"net/http"
	"time"

	"github.com/FunnyFoXD/go_final_project/helpers"
)

func NewDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	if nowStr == "" || dateStr == "" || repeat == "" {
		http.Error(w, "missing some parameters", http.StatusBadRequest)
		return
	}

	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		http.Error(w, "Invalid now parameter", http.StatusBadRequest)
		return
	}

	_, err = time.Parse("20060102", dateStr)
	if err != nil {
		http.Error(w, "Invalid date parameter", http.StatusBadRequest)
		return
	}

	nextDate, err := helpers.NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if nextDate == "" {
		http.Error(w, "No next date", http.StatusBadRequest)
		return
	}

	w.Write([]byte(nextDate))
}