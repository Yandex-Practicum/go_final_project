package main

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

	search := r.URL.Query().Get("search")
	var tasks []Task
	var query string
	var args []interface{}

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

	for rows.Next() {
		var task Task
		var id int
		err := rows.Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"Ошибка чтения БД: %s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		task.ID = strconv.Itoa(id)
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"Ошибка обработки: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"tasks": tasks,
	}

	if len(tasks) == 0 {
		response["tasks"] = []Task{}
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// ПОЛУЧЕНИЕ ЗАДАЧИ по ID
func getTaskHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"error":"Не указан id"}`, http.StatusBadRequest)
		return
	}

	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?"
	row := db.QueryRow(query, id)
	var task Task
	var taskID int
	err := row.Scan(&taskID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err == sql.ErrNoRows {
		http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"Ошибка поиска: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	task.ID = strconv.Itoa(taskID)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}
