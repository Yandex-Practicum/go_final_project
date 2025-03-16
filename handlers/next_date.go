package handlers

import (
	"github.com/qavaleria/go_final_project/tasks"
	"net/http"
	"time"
)

// HandleNextDate обработчик для API запроса /api/nextdate
func HandleNextDate(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	now, err := time.Parse(tasks.FormatDate, nowStr)
	if err != nil {
		http.Error(w, "Invalid now format", http.StatusBadRequest)
		return
	}

	nextDate, err := tasks.NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(nextDate))
}
