package tasks

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

	var date, repeat string
	err := db.QueryRow("SELECT date, repeat FROM scheduler WHERE id = ?", id).Scan(&date, &repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf(`{"error":"Ошибка поиска: %s"}`, err.Error()), http.StatusInternalServerError)
		}
		return
	}

	if repeat == "" {
		if _, err := db.Exec("DELETE FROM scheduler WHERE id = ?", id); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"Ошибка удаления: %s"}`, err.Error()), http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{}`)
		}
		return
	}

	nextDate, err := NextDate(time.Now(), date, repeat)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"Ошибка расчета даты: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	if _, err := db.Exec("UPDATE scheduler SET date = ? WHERE id = ?", nextDate, id); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"Ошибка обновления даты: %s"}`, err.Error()), http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{}`)
	}
}
