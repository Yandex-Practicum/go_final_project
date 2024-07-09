package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func NextDateGETHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.URL.Query().Get("now")
	date := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	if nowStr == "" || date == "" || repeat == "" {
		http.Error(w, "нет параметров", http.StatusBadRequest)
		return
	}

	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		http.Error(w, "неверный формат времени", http.StatusBadRequest)
		return
	}

	nextDate, err := NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, nextDate)
	//w.Write([]byte(nextDate))
}

func NextDate(now time.Time, date string, repeat string) (string, error) {
	// Парсим дату начала
	startDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", errors.New("неверный формат даты")
	}

	// Если правило пустое - ошибка
	if repeat == "" {
		return "", errors.New("правило не может быть пустым")
	}

	// Определяем следующую дату на основе правила
	switch {
	case strings.HasPrefix(repeat, "d "):
		// Обработка правила "d <число>"
		daysStr := strings.TrimPrefix(repeat, "d ")
		daysToAdd, err := strconv.Atoi(daysStr)
		if err != nil || daysToAdd <= 0 || daysToAdd > 400 {
			return "", errors.New("неверное правило для дней")
		}

		endDate := startDate.AddDate(0, 0, daysToAdd)
		for endDate.Before(now) {
			endDate = endDate.AddDate(0, 0, daysToAdd)
		}
		return endDate.Format("20060102"), nil

	case repeat == "y":
		// Обработка правила "y"
		endDate := startDate.AddDate(1, 0, 0)
		for endDate.Before(now) {
			endDate = endDate.AddDate(1, 0, 0)
		}

		return endDate.Format("20060102"), nil

	default:
		return "", errors.New("неверное правило")
	}
}
