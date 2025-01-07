package handlers

import (
	"net/http"
	"time"

	"github.com/FunnyFoXD/go_final_project/helpers"
)

// NextDateHandler is a handler for "/api/nextdate" endpoint.
// It takes "now", "date" and "repeat" parameters and returns the next date
// based on the given parameters.
// The "now" parameter should be given in the format "YYYYMMDD".
// The "date" parameter should be given in the same format.
// The "repeat" parameter can be "y" for yearly repeat or "d <days>" for daily repeat.
// If the parameters are invalid or if the date is invalid, it returns an error with HTTP status code 400.
// If the date is successfully calculated, it returns the next date in the same format with HTTP status code 200.
func NextDateHandler(w http.ResponseWriter, r *http.Request) {
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

	w.Write([]byte(nextDate))
}