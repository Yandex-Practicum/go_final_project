package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

var (
	dateRegex   = regexp.MustCompile(`^\d{8}$`)
	repeatRegex = regexp.MustCompile(`^(d \d+|y|w (\d,?)+)$`)
)

func handleCreateTask(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	const dateFormat = "20060102"
	var task Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, `{"error":"Ошибка десериализации JSON"}`, http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		http.Error(w, `{"error":"Не указан заголовок задачи"}`, http.StatusBadRequest)
		return
	}
	if task.Date == "" {
		task.Date = time.Now().Format(dateFormat)
	}

	if !dateRegex.MatchString(task.Date) {
		http.Error(w, `{"error":"Дата указана в неверном формате"}`, http.StatusBadRequest)
		return
	}

	date, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		http.Error(w, `{"error":"Дата указана в неверном формате"}`, http.StatusBadRequest)
		return
	}

	now := time.Now()
	if date.Format(dateFormat) < now.Format(dateFormat) {
		if task.Repeat == "" {
			date = now
		} else {
			nextDateStr, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				http.Error(w, `{"error":"Ошибка при вычислении следующей даты"}`, http.StatusBadRequest)
				return
			}
			date, err = time.Parse(dateFormat, nextDateStr)
			if err != nil {
				http.Error(w, `{"error":"Ошибка обработки следующей даты"}`, http.StatusInternalServerError)
				return
			}
		}
	}
	task.Date = date.Format(dateFormat)

	if task.Repeat != "" && !repeatRegex.MatchString(task.Repeat) {
		http.Error(w, `{"error":"Неверный формат правила повторения"}`, http.StatusBadRequest)
		return
	}

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		http.Error(w, `{"error":"Ошибка базы данных"}`, http.StatusInternalServerError)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		http.Error(w, `{"error":"Не удалось получить ID"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"id": strconv.FormatInt(id, 10)})
}
func handleGetTasks(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	const limit = 50
	query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT ?`

	rows, err := db.Query(query, limit)
	if err != nil {
		http.Error(w, `{"error":"Ошибка базы данных"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tasks []map[string]string
	for rows.Next() {
		var id int
		var date, title, comment, repeat string
		if err := rows.Scan(&id, &date, &title, &comment, &repeat); err != nil {
			http.Error(w, `{"error":"Ошибка чтения данных"}`, http.StatusInternalServerError)
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

	if tasks == nil {
		tasks = []map[string]string{}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"tasks": tasks})
}
func handleGetTaskByID(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"error":"Не указан идентификатор"}`, http.StatusBadRequest)
		return
	}

	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	row := db.QueryRow(query, id)

	var task Task
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err == sql.ErrNoRows {
		http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, `{"error":"Ошибка базы данных"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(task)
}
func handleUpdateTask(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, `{"error":"Ошибка десериализации JSON"}`, http.StatusBadRequest)
		return
	}

	if task.ID == "" {
		http.Error(w, `{"error":"Не указан ID задачи"}`, http.StatusBadRequest)
		return
	}
	if task.Title == "" {
		http.Error(w, `{"error":"Не указан заголовок задачи"}`, http.StatusBadRequest)
		return
	}
	if !dateRegex.MatchString(task.Date) {
		http.Error(w, `{"error":"Дата указана в неверном формате"}`, http.StatusBadRequest)
		return
	}

	// Проверяем, можно ли распарсить дату
	if _, err := time.Parse("20060102", task.Date); err != nil {
		http.Error(w, `{"error":"Дата указана в неверном формате"}`, http.StatusBadRequest)
		return
	}

	if task.Repeat != "" && !repeatRegex.MatchString(task.Repeat) {
		http.Error(w, `{"error":"Неверный формат правила повторения"}`, http.StatusBadRequest)
		return
	}

	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		http.Error(w, `{"error":"Ошибка базы данных"}`, http.StatusInternalServerError)
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		http.Error(w, `{"error":"Ошибка обработки результата"}`, http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{})
}
func handleTaskDone(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"error":"Не указан идентификатор"}`, http.StatusBadRequest)
		return
	}

	var task Task
	query := `SELECT id, date, repeat FROM scheduler WHERE id = ?`
	err := db.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Repeat)
	if err == sql.ErrNoRows {
		http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, `{"error":"Ошибка базы данных"}`, http.StatusInternalServerError)
		return
	}

	if task.Repeat == "" {
		// Если repeat пуст, удаляем задачу
		_, err := db.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
		if err != nil {
			http.Error(w, `{"error":"Ошибка удаления задачи"}`, http.StatusInternalServerError)
			return
		}
	} else {
		// Вычисляем следующую дату
		now := time.Now()
		nextDate, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			http.Error(w, `{"error":"Ошибка вычисления следующей даты"}`, http.StatusBadRequest)
			return
		}

		// Обновляем дату задачи
		_, err = db.Exec(`UPDATE scheduler SET date = ? WHERE id = ?`, nextDate, id)
		if err != nil {
			http.Error(w, `{"error":"Ошибка обновления даты"}`, http.StatusInternalServerError)
			return
		}
	}

	json.NewEncoder(w).Encode(map[string]string{})
}
func handleDeleteTask(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"error":"Не указан идентификатор задачи"}`, http.StatusBadRequest)
		return
	}

	query := `DELETE FROM scheduler WHERE id = ?`
	res, err := db.Exec(query, id)
	if err != nil {
		http.Error(w, `{"error":"Ошибка базы данных"}`, http.StatusInternalServerError)
		return
	}

	affected, err := res.RowsAffected()
	if err != nil {
		http.Error(w, `{"error":"Ошибка при проверке удаления"}`, http.StatusInternalServerError)
		return
	}

	if affected == 0 {
		http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{})
}
