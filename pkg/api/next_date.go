package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func afterNow(date, now time.Time) bool {

	now = now.Truncate(24 * time.Hour)
	date = date.Truncate(24 * time.Hour)
	return date.After(now)

}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {

	// Получаем date
	date, err := time.Parse("20060102", dstart)

	if err != nil {
		return "", fmt.Errorf("cannot to parse dstart: %w", err)
	}

	// Проверяем, что правило не пустое
	if repeat == "" {
		return "", fmt.Errorf("repeat rule cannot be empty")
	}

	repeatParts := strings.Split(repeat, " ")

	switch repeatParts[0] {
	case "d":
		if len(repeatParts) != 2 {
			return "", fmt.Errorf("invalid format: %w", err)
		}

		interval, err := strconv.Atoi(repeatParts[1])
		if err != nil {
			return "", fmt.Errorf("invalid interval format. Cannot convert: %w", err)
		}

		if interval < 1 || interval > 400 {
			return "", fmt.Errorf("interval should be between 1 and 400")
		}

		for {
			date = date.AddDate(0, 0, interval)
			if afterNow(date, now) {
				return date.Format("20060102"), nil
			}
		}

	case "y":
		if len(repeatParts) != 1 {
			return "", fmt.Errorf("invalid year format")
		}

		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				return date.Format("20060102"), nil
			}
		}

	default:
		return "", fmt.Errorf("unknown format: %s", repeatParts[0])
	}

}

func NextDateHandler(w http.ResponseWriter, r *http.Request) {

	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")
	nowStr := r.FormValue("now")

	var now time.Time

	if nowStr == "" {

		now = time.Now().UTC().Truncate(24 * time.Hour)

	} else {
		var err error
		now, err = time.Parse("20060102", nowStr)
		if err != nil {
			writeJsonError(w, "incorrect parameter 'now': "+err.Error(), http.StatusBadRequest)
			return
		}
		now = now.Truncate(24 * time.Hour)
	}

	nextDate, err := NextDate(now, dateStr, repeat)
	if err != nil {
		writeJsonError(w, "cannot calculate next date: "+err.Error(), http.StatusBadRequest)
		return
	}

	writeJson(w, map[string]string{"date": nextDate})

}
