package handlers

import (
	"encoding/json"
	"net/http"
	"pwd/services"
	"time"
)

func NextDateHandler(w http.ResponseWriter, r *http.Request) {

	nowStr := r.URL.Query().Get("now")
	dateStr := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	if nowStr == "" || dateStr == "" || repeat == "" {
		http.Error(w, "недостаточно значений", http.StatusBadRequest)
		return
	}

	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		http.Error(w, "недопустимый формат даты", http.StatusBadRequest)
		return
	}

	nextDate, err := services.NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := map[string]string{"next_date": nextDate}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
