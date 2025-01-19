package handlers

import (
	"log"
	"net/http"
	"time"

	"final-project-bronner/go/nextdate"
)

const DateFormat = "20060102"

func (h *Handler) NextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	now, err := time.Parse(DateFormat, nowStr)
	if err != nil {
		http.Error(w, "Невалидный формат даты, тебуемый формат: YYYYMMDD", http.StatusBadRequest)
		return
	}

	nextDate, err := nextdate.NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, err := w.Write([]byte(nextDate)); err != nil {
		log.Printf("Запрос завершился с ошибкой: %s", err)
	}
}
