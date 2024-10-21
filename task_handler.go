package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	//"log"
	"net/http"
	"strconv"
	"time"

	_ "modernc.org/sqlite"
)


const limit int8 = 50

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func Check(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		AddTask(w, r)
	case http.MethodGet:
		GetTask(w, r)
	case http.MethodPut:
		EditTask(w, r)
	case http.MethodDelete:
		DeleteTask(w, r)
	default:
		http.Error(w, `{"error":"Invalid request method"}`, http.StatusMethodNotAllowed)
	}
}

func AddTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var task Task
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&task); err != nil {
		http.Error(w, `{"error":"Invalid request payload"}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if task.Title == "" {
		http.Error(w, `{"error":"Task title is required"}`, http.StatusBadRequest)
		return
	}

	if task.Date == "" {
		task.Date = time.Now().Format(formatDate)
	} else {
		_, err := time.Parse(formatDate, task.Date)
		if err != nil {
			http.Error(w, `{"error":"Invalid date format, should be YYYYMMDD"}`, http.StatusBadRequest)
			return
		}
	}

	now := time.Now().Format(formatDate)
	if task.Date < now && task.Repeat == "" {
		task.Date = now
	}

	if task.Repeat != "" {
		parsedDate, err := time.Parse(formatDate, task.Date)
		if err != nil {
			http.Error(w, `{"error":"Invalid date format, should be YYYYMMDD"}`, http.StatusBadRequest)
			return
		}

		nextDate, err := NextDate(parsedDate, task.Date, task.Repeat)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
			return
		}

		if nextDate < now {
			nextDate = now
		}

		task.Date = nextDate
	}

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		http.Error(w, `{"error":"Failed to add task to the database"}`, http.StatusInternalServerError)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		http.Error(w, `{"error":"Failed to retrieve task ID"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"id":"%d"}`, id)

}

func GetTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"Invalid request method"}`, http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"error":"Task ID is required"}`, http.StatusBadRequest)
		return
	}

	err := DB.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id=?", id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err == sql.ErrNoRows {
		http.Error(w, `{"error":"Task is missing"}`, http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, `{"error":"Error receiving  task"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		http.Error(w, `{"error":"Ошибка записи ответа"}`, http.StatusInternalServerError)
		return
	}

}

func GetTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"Invalid request method"}`, http.StatusMethodNotAllowed)
		return
	}

	search := r.URL.Query().Get("search")
	query := `SELECT id, date, title, comment, repeat FROM scheduler`
	args := []interface{}{}

	if search != "" {

		if parsedDate, err := time.Parse("02.01.2006", search); err == nil {
			query += ` WHERE date = ?`
			args = append(args, parsedDate.Format(formatDate))
		} else {

			query += ` WHERE title LIKE ? OR comment LIKE ?`
			searchTerm := "%" + search + "%"
			args = append(args, searchTerm, searchTerm)
		}
	}

	query += ` ORDER BY date LIMIT ?`
	args = append(args, limit)

	rows, err := DB.Query(query, args...)
	if err != nil {
		http.Error(w, `{"error":"Failed to retrieve tasks from the database"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	tasks := []map[string]string{}
	for rows.Next() {
		var id int
		var date, title, comment, repeat string
		if err := rows.Scan(&id, &date, &title, &comment, &repeat); err != nil {
			http.Error(w, `{"error":"Failed to scan task from the database"}`, http.StatusInternalServerError)
			return
		}

		task := map[string]string{
			"id":      strconv.Itoa(id),
			"date":    date,
			"title":   title,
			"comment": comment,
			"repeat":  repeat,
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, `{"error":"Database error occurred"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"tasks": tasks,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, `{"error":"Failed to encode tasks to JSON"}`, http.StatusInternalServerError)
		return
	}

}

func EditTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&task); err != nil {
		http.Error(w, `{"error":"Invalid request payload"}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if task.ID == "" {
		http.Error(w, `{"error":"Task ID is required"}`, http.StatusBadRequest)
		return
	}

	trueTask, err := getID(task.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusNotFound)
		return
	}

	if task.Title == "" {
		http.Error(w, `{"error":"Task title is required"}`, http.StatusBadRequest)
		return
	}

	if task.Date == "" {
		task.Date = trueTask.Date
	} else {
		_, err := time.Parse(formatDate, task.Date)
		if err != nil {
			http.Error(w, `{"error":"Invalid date format, should be YYYYMMDD"}`, http.StatusBadRequest)
			return
		}
	}

	if task.Repeat != "" {
		nextDate, err := NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
			return
		}
		task.Date = nextDate
	}

	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	_, err = DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		http.Error(w, `{"error":"Failed to update task in the database"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(`{}`)); err != nil {
		fmt.Printf("Error writing response: %v", err)
	}
}

func getID(id string) (*Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	row := DB.QueryRow(query, id)

	var task Task
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Task not found")
		}
		return nil, err
	}

	return &task, nil

}

func DoneTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Invalid request method"}`, http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"error":"Task ID is required"}`, http.StatusBadRequest)
		return
	}

	now := time.Now()

	var task Task
	err := DB.QueryRow("SELECT date, repeat FROM scheduler WHERE id = ?", id).Scan(&task.Date, &task.Repeat)
	if err != nil {
		http.Error(w, `{"error":"Task not found"}`, http.StatusNotFound)
		return
	}

	if task.Repeat == "" {
		if _, err = DB.Exec("DELETE FROM scheduler WHERE id = ?", id); err != nil {
			http.Error(w, `{"error":"Failed to delete task"}`, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{}`)
		return
	}

	nextDate, err := NextDate(now, task.Date, task.Repeat)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	if _, err := DB.Exec("UPDATE scheduler SET date = ? WHERE id = ?", nextDate, id); err != nil {
		http.Error(w, `{"error":"Failed to update task date"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{}`)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, `{"error":"Invalid request method"}`, http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"error":"Task ID is required"}`, http.StatusBadRequest)
		return
	}

	_, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, `{"error":"Invalid task ID"}`, http.StatusBadRequest)
		return
	}

	query := `DELETE FROM scheduler WHERE id = ?`
	result, err := DB.Exec(query, id)
	if err != nil {
		http.Error(w, `{"error":"Failed to delete task from the database"}`, http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, `{"error":"Failed to retrieve affected rows"}`, http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, `{"error":"Task not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "{}")
}
