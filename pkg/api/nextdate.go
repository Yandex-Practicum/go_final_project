package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	date, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", fmt.Errorf("неверный формат даты: %w", err)
	}

	if repeat == "" {
		return "", fmt.Errorf("правило повторения не указано")
	}

	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", fmt.Errorf("пустое правило повторения")
	}

	nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", fmt.Errorf("неверный формат правила d")
		}

		interval, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", fmt.Errorf("интервал должен быть числом")
		}

		if interval < 1 || interval > 400 {
			return "", fmt.Errorf("интервал должен быть от 1 до 400")
		}

		for !date.After(nowDate) {
			date = date.AddDate(0, 0, interval)
		}

		// Если начальная дата уже находится в будущем,
		// возвращаем следующее повторение.
		if date.Format("20060102") == dstart {
			date = date.AddDate(0, 0, interval)
		}

	case "y":
		if len(parts) != 1 {
			return "", fmt.Errorf("неверный формат правила y")
		}

		for !date.After(nowDate) {
			date = date.AddDate(1, 0, 0)
		}

		if date.Format("20060102") == dstart {
			date = date.AddDate(1, 0, 0)
		}

	default:
		return "", fmt.Errorf("неизвестное правило повторения")
	}

	return date.Format("20060102"), nil
}

func sameDay(date, now time.Time) bool {
	return date.Year() == now.Year() &&
		date.Month() == now.Month() &&
		date.Day() == now.Day()
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	now := r.URL.Query().Get("now")
	date := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	var nowTime time.Time
	var err error

	if now == "" {
		nowTime = time.Now()
	} else {
		nowTime, err = time.Parse("20060102", now)
		if err != nil {
			http.Error(w, "неверный формат даты now", http.StatusBadRequest)
			return
		}
	}

	nextDate, err := NextDate(nowTime, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, _ = w.Write([]byte(nextDate))
}
