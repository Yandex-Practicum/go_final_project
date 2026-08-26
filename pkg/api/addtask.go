package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/UTKIN-AI/go_final_project/pkg/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	if task.Title == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "title is required",
		})
		return
	}

	// Проверяем дату и правило повторения
	err = checkDate(&task)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	// добавляем задачу в БД.
	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Возвращаем ID созданной задачи.
	err = writeJSON(w, http.StatusOK, map[string]string{
		"id": strconv.FormatInt(id, 10),
	})
	if err != nil {
		return
	}

}

func checkDate(task *db.Task) error {
	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format(dateLayout)
	}

	t, err := time.Parse(dateLayout, task.Date)
	if err != nil {
		return fmt.Errorf("invalid date format: %w", err)
	}

	today, _ := time.Parse(dateLayout, now.Format(dateLayout))

	if t.Before(today) {
		if task.Repeat == "" {
			task.Date = today.Format(dateLayout)
		} else {
			next, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return fmt.Errorf("invalid repeat rule: %w", err)
			}

			task.Date = next
		}
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}
