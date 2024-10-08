package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"
)

// ПОМЕТКА О ВЫПОЛЕНИИ ЗАДАЧИ
func taskDoneHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"error":"Не указан id"}`, http.StatusBadRequest)
		return
	}

	query := "SELECT date, repeat FROM scheduler WHERE id = ?"
	row := db.QueryRow(query, id)

	var date string
	var repeat string
	err := row.Scan(&date, &repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf(`{"error":"Ошибка поиска: %s"}`, err.Error()), http.StatusInternalServerError)
		}
		return
	}

	if repeat == "" {
		_, err = db.Exec("DELETE FROM scheduler WHERE id = ?", id)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"Ошибка удаления: %s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{}`)
		return
	}

	now := time.Now()
	nextDate, err := NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"Ошибка расчета даты: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	query = `UPDATE scheduler SET date = ? WHERE id = ?`
	_, err = db.Exec(query, nextDate, id)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"Ошибка обновления даты: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{}`)
}
