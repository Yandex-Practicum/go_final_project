package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go_final_project/pkg/db"
)

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(v)
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var t db.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	if t.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": "title is required"})
		return
	}

	// нормализуем дату «сегодня» к полуночи
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayStr := today.Format(isoDate)

	// если дата не указана, ставим сегодня
	if t.Date == "" {
		t.Date = todayStr
	} else {
		parsed, err := time.Parse(isoDate, t.Date)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			writeJSON(w, map[string]string{"error": "invalid date format"})
			return
		}
		// корректируем только если дата в прошлом
		if parsed.Before(today) {
			if t.Repeat == "" {
				t.Date = todayStr
			} else {
				next, err := NextDate(today, t.Date, t.Repeat)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					writeJSON(w, map[string]string{"error": fmt.Sprintf("repeat rule error: %v", err)})
					return
				}
				t.Date = next
			}
		}
	}

	id, err := db.AddTask(&t)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, map[string]string{"id": fmt.Sprint(id)})
}
