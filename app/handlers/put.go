package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"go_final/app/database"
	"go_final/app/models"
	"go_final/app/tasks"
)

func PutTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var task models.Remind
	if err := decoder.Decode(&task); err != nil {
		http.Error(w, `{"error":"Error JSON decoding"}`, http.StatusBadRequest)
		return
	}

	if task.ID == "" {
		http.Error(w, `{"error":"Task id not specified"}`, http.StatusBadRequest)
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
		parsedDate, err := time.Parse("20060102", task.Date)
		if err != nil {
			http.Error(w, `{"error":"Incorrect date format"}`, http.StatusBadRequest)
			return
		}
		if parsedDate.Before(time.Now()) {
			if task.Repeat == "" {
				task.Repeat = today
			} else {
				nextDate, err := tasks.NextDate(now, task.Date, task.Repeat)
				if err != nil {
					http.Error(w, `{"error":"Incorrect task repeat rule"}`, http.StatusBadRequest)
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
	err := database.UpdateTask(task)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error":"Task not found"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"error":"Task update error"}`, http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{}`))
}
