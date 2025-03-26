package handler

import (
	"fmt"
	"net/http"
	"time"

	nextdate "github.com/askoren1/go_final_project/internal/next_date"
)

// Функция NextDate для вычисления следующей даты выполнения задачи на основе заданной даты и правила повторения
func (h *Handler) NextDate(w http.ResponseWriter, r *http.Request) {

	// Получение данных из запроса
	nowTimeStr := r.FormValue("now")
	nowTime, _ := time.Parse(Layout, nowTimeStr)	
	dateStr := r.FormValue("date")
	repeatStr := r.FormValue("repeat")

	if dateStr == "" {
		dateStr = time.Now().Format(Layout)
	}

	// Вычисление следующей даты
	date, err := nextdate.NextDate(nowTime, dateStr, repeatStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintln(w, date)
}
