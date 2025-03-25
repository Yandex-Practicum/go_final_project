package scheduler

import (
	"fmt"
	"net/http"
	"time"
)

// NexDateHandler обрабатывает запросы к  /api/nextdate
func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем парметры из запроса
	nowStr := r.URL.Query().Get("now")
	dateStr := r.URL.Query().Get("date")
	repeatRule := r.URL.Query().Get("repeat")

	// Паррсим текущще время
	now, err := parseTime(nowStr)
	if err != nil {
		http.Error(w, "Неверный формат даты1", http.StatusBadRequest)
		return
	}

	//Вычисляем следующую дату
	NextDate, err := NextDate(now, dateStr, repeatRule)
	if err != nil {
		http.Error(w, "Неверный формат даты2", http.StatusBadRequest)
		return
	}

	//Возвращаем результат
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, NextDate)
}

// parseTime парсит время
func parseTime(nowStr string) (time.Time, error) {
    if nowStr == "" {
        return time.Now().UTC(), nil
    }
    timeParse, err := time.Parse("20060102", nowStr)
    return timeParse, err
}
