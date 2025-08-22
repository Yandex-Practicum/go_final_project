package api

import (
	"fmt"
	"go_final_project/pkg/nextdate"
	"log"
	"net/http"
	"time"
)

// Обработчик для вычисления следующей даты
func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	// Извлекаем параметры из запроса
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeatStr := r.FormValue("repeat")

	// Логируем параметры для отладки
	log.Printf("Received parameters: now=%s, date=%s, repeat=%s", nowStr, dateStr, repeatStr)

	// Преобразуем nowStr в формат time
	now, err := time.Parse(nextdate.TimeFormat, nowStr)
	if err != nil {
		http.Error(w, "Ошибка парсинга now", http.StatusBadRequest)
		return
	}

	// Преобразуем dateStr в формат time
	date, err := time.Parse(nextdate.TimeFormat, dateStr)
	if err != nil {
		http.Error(w, "Ошибка парсинга date", http.StatusBadRequest)
		return
	}

	nextDate, err := nextdate.NextDate(now, date.Format(nextdate.TimeFormat), repeatStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка вычисления следующей даты: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("Next date: %s", nextDate)
	fmt.Fprintf(w, "%s", nextDate)
}
