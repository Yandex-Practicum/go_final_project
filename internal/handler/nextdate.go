package handler

import (
	"go_final_project/internal/service"
	"net/http"
	"time"
)

type NextDateHandler struct{}

func NewNextDateHandler() *NextDateHandler {
	return &NextDateHandler{}
}

func (h *NextDateHandler) HandleNextDate(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		http.Error(w, "Неверный текущий формат даты", http.StatusBadRequest)
		return
	}

	nextDate, err := service.NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(nextDate))
}
