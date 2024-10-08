package tasks

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// ПОЛУЧЕНИЕ всех ЗАДАЧ
func tasksHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	id := r.URL.Query().Get("id")
	search := r.URL.Query().Get("search")

	var query string
	var args []interface{}

	if id != "" {
		query = "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?"
		args = append(args, id)
		var task Task
		if err := db.QueryRow(query, args...).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err == sql.ErrNoRows {
			http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"Ошибка поиска: %s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(task)
		return
	}

	// Если ID не указан, ищем список задач с возможностью поиска
	if search != "" {
		if searchDate, err := time.Parse("02.01.2006", search); err == nil {
			query = "SELECT * FROM scheduler WHERE date = ? ORDER BY date LIMIT 50"
			args = append(args, searchDate.Format(dateFormat))
		} else {
			query = "SELECT * FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT 50"
			searchPattern := "%" + search + "%"
			args = append(args, searchPattern, searchPattern)
		}
	} else {
		query = "SELECT * FROM scheduler ORDER BY date LIMIT 50"
	}
	rows, err := db.Query(query, args...)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"Ошибка чтения БД: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	tasks := []Task{}
	for rows.Next() {
		var task Task
		var id int
		if err := rows.Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"Ошибка чтения БД: %s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		task.ID = strconv.Itoa(id)
		tasks = append(tasks, task)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"tasks": tasks})
}
