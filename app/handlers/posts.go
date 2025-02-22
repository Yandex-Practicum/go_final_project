package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"go_final/app/database"
	"go_final/app/models"
	"go_final/app/tasks"
)

func PostTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var task = models.Remind{}
	if err := decoder.Decode(&task); err != nil {
		http.Error(w, `{"error":"Error JSON encoding"}`, http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		http.Error(w, `{"error":"Task title not specified"}`, http.StatusBadRequest)
		return
	}

	now := time.Now()
	today := now.Format("20060102")

	if task.Date == "" || task.Date == "today" {
		task.Date = today
	} else {
		// парсим дату задачи
		parsedDate, err := time.Parse("20060102", task.Date)
		if err != nil {
			http.Error(w, `{"error":"Incorrect date format"}`, http.StatusBadRequest)
			return
		}
		if parsedDate.Before(time.Now()) {
			if task.Repeat == "" {
				task.Date = today
			} else {
				nextDate, err := tasks.NextDate(time.Now(), task.Date, task.Repeat)
				if err != nil {
					http.Error(w, `{"error":"Incorrect task repeat rule}`, http.StatusBadRequest)
					return
				}
				if task.Date != today {
					task.Date = nextDate
				}
			}
		}
	}
	if task.Repeat != "" {
		_, err := tasks.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			http.Error(w, `{"error":"Incorrect task repeat rule"}`, http.StatusBadRequest)
			return
		}
	}
	id, err := database.InsertIntoDB(task)
	idForResponse := strconv.Itoa(int(id))

	if err != nil {
		http.Error(w, `{"error":"Error access to database"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(map[string]string{
		"id": idForResponse,
	})
}
