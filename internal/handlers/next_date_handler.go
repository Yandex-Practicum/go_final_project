package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func nextDate(now time.Time, dateStr, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("необходимо указать интервал повторения (repeat)")
	}

	date, err := time.Parse("20060102", dateStr)
	if err != nil {
		return "", fmt.Errorf("неверный формат даты, не удалось преобразовать в корректную дату")
	}

	if now.IsZero() {
		now = time.Now()
	}

	parts := strings.Split(repeat, " ")
	if len(parts) < 1 || len(parts) > 2 {
		return "", fmt.Errorf("неверный формат repeat, требуется указать интервал и (при необходимости) количество")
	}

	interval := parts[0]
	var next time.Time

	switch interval {
	case "y":
		next = date.AddDate(1, 0, 0)
		for next.Before(now) {
			next = next.AddDate(1, 0, 0)
		}

	case "d":
		if len(parts) > 1 {
			days, err := strconv.Atoi(parts[1])
			if err != nil || days <= 0 || days > 400 {
				return "", fmt.Errorf("некорректный аргумент для 'd': %v, должно быть числом от 1 до 400", parts[1])
			}
			next = date.AddDate(0, 0, days)
			for next.Before(now) {
				next = next.AddDate(0, 0, days)
			}
		} else {
			return "", fmt.Errorf("не указан интервал для дней после 'd'")
		}

	default:
		return "", fmt.Errorf("неподдерживаемый интервал: %s", interval)
	}

	if next.IsZero() {
		return "", fmt.Errorf("не удалось вычислить следующую дату")
	}

	return next.Format("20060102"), nil
}

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")
	nowStr := r.URL.Query().Get("now")

	if dateStr == "" || repeat == "" || nowStr == "" {
		http.Error(w, "Параметры 'date', 'repeat' и 'now' обязательны", http.StatusBadRequest)
		return
	}

	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		http.Error(w, "Неверный формат 'now'", http.StatusBadRequest)
		return
	}

	next, err := nextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(next))
}
