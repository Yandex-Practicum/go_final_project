package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/FunnyFoXD/go_final_project/databases"
	"github.com/FunnyFoXD/go_final_project/helpers"
	"github.com/FunnyFoXD/go_final_project/models"
)

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var updatedTask models.TaskFromDB
	var buf bytes.Buffer
	var now = time.Now().Truncate(24 * time.Hour)

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, `{"error":"can't read body"}`, http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &updatedTask); err != nil {
		http.Error(w, `{"error":"can't unmarshal body"}`, http.StatusBadRequest)
		return
	}

	if updatedTask.Title == "" {
		http.Error(w, `{"error":"title is empty"}`, http.StatusBadRequest)
		return
	}

	if updatedTask.Date == "" || updatedTask.Date == "today" || updatedTask.Date == "Today" {
		updatedTask.Date = now.Format("20060102")
	}

	taskParse, err := time.Parse("20060102", updatedTask.Date)
	if err != nil {
		http.Error(w, `{"error":"invalid date format"}`, http.StatusBadRequest)
		return
	}

	if taskParse.Before(now) {
		if updatedTask.Repeat == "" {
			updatedTask.Date = now.Format("20060102")
		} else {
			nextDate, err := helpers.NextDate(now, updatedTask.Date, updatedTask.Repeat)
			if err != nil {
				http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusBadRequest)
				return
			}

			updatedTask.Date = nextDate
		}
	}

	err = databases.UpdateTaskByID(updatedTask)
	if err == sql.ErrNoRows {
		http.Error(w, `{"error":"task not found"}`, http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte("{}"))
}
