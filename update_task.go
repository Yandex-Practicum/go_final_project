package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// ОБНОВЛЕНИЕ ЗАДАЧИ
func editTaskHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var task Task
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&task)
	if err != nil {
		http.Error(w, `{"error":"Ошибка декодирования JSON"}`, http.StatusBadRequest)
		return
	}
	if task.ID == "" {
		http.Error(w, `{"error":"Не указан id"}`, http.StatusBadRequest)
		return
	}
	if task.Title == "" {
		http.Error(w, `{"error":"Не указан заголовок"}`, http.StatusBadRequest)
		return
	}
	if task.Date == "" {
		task.Date = time.Now().Format(dateFormat)
	} else {
		_, err := time.Parse(dateFormat, task.Date)
		if err != nil {
			http.Error(w, `{"error":"Неверный формат даты"}`, http.StatusBadRequest)
			return
		}
	}

	now := time.Now()
	tDate, _ := time.Parse(dateFormat, task.Date)
	if tDate.Before(now) {
		if task.Repeat == "" {
			task.Date = now.Format(dateFormat)
		} else {
			nextDate, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
				return
			}
			task.Date = nextDate
		}
	}
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"Ошибка обновления: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"Ошибка чтения из БД: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{}`)
}
