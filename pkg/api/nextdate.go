package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const layout = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("empty repeat rule")
	}

	startDate, err := time.Parse(layout, dstart)
	if err != nil {
		return "", fmt.Errorf("invalid dstart format: %w", err)
	}

	parts := strings.Split(repeat, " ")
	rule := parts[0]

	switch rule {
	case "d":
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid daily rule format")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days <= 0 || days > 400 {
			return "", fmt.Errorf("invalid number of days")
		}

		for {
			startDate = startDate.AddDate(0, 0, days)
			if startDate.After(now) {
				break
			}
		}

	case "y":
		for {
			startDate = startDate.AddDate(1, 0, 0)
			if startDate.After(now) {
				break
			}
		}

	default:
		return "", fmt.Errorf("unsupported repeat rule")
	}

	return startDate.Format(layout), nil
}

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	layout := "20060102"

	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	var err error

	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(layout, nowStr)
		if err != nil {
			http.Error(w, "invalid now param", http.StatusBadRequest)
			return
		}
	}

	next, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(next))
}
