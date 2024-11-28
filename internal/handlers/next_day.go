package handlers

import (
	"fmt"
	"my_education/go/go_final_project/internal/logic" // Путь к вашему пакету с функцией NextDate
	"net/http"
	"time"
)

// NextDateHandler обрабатывает запросы на вычисление следующей даты.
func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	// Парсим дату now из строки
	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		http.Error(w, "некорректная дата now", http.StatusBadRequest)
		return
	}

	// Вызываем функцию для получения следующей даты
	nextDate, err := logic.NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, fmt.Sprintf("ошибка: %v", err), http.StatusBadRequest)
		return
	}

	// Возвращаем результат в формате текста
	fmt.Fprintln(w, nextDate)
}
