package parsedate

import (
	"fmt"
	"net/http"
	"time"
)

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.URL.Query().Get("now")
	dateStr := r.URL.Query().Get("date")
	repeatStr := r.URL.Query().Get("repeat")

	now, err := parseTime(nowStr)
	if err != nil {
		http.Error(w, "bad format", http.StatusBadRequest)
		return
	}

	NextDate, err := NextDate(now, dateStr, repeatStr)
	if err != nil {
		http.Error(w, "bad format", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, NextDate)
}

func parseTime(nowStr string) (time.Time, error) {
	timeParse, err := time.Parse("20060102", nowStr)
	return timeParse, err
}
