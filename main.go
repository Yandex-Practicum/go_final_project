package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3" // Импортируем драйвер sqlite3
)

const dateFormat = "20060102"

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

var db *sql.DB

func initDB() error {
	// Получаем текущий рабочий каталог
	appPath, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %v", err)
	}

	// Определяем полный путь к файлу базы данных
	dbFile := filepath.Join(appPath, "scheduler.db")

	// Открываем или создаем базу данных
	db, err = sql.Open("sqlite3", dbFile)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	// Проверяем, существует ли таблица scheduler
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS scheduler (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            date TEXT NOT NULL CHECK (LENGTH(date) <= 8),
            title TEXT NOT NULL CHECK (LENGTH(title) <= 255),
            comment TEXT CHECK (LENGTH(comment) <= 1024),
            repeat TEXT CHECK (LENGTH(repeat) <= 128)
        );
        CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);
    `)
	if err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}

	log.Println("Database initialized successfully.")
	return nil
}

func NextDate(now time.Time, date string, repeat string) (string, error) {
	d, err := time.Parse(dateFormat, date)
	if err != nil {
		return "", fmt.Errorf("invalid date format: %v", err)
	}

	switch {
	case repeat == "":
		return "", fmt.Errorf("repeat rule is empty")
	case repeat == "y":
		d = d.AddDate(1, 0, 0)
		for d.Format(dateFormat) <= now.Format(dateFormat) {
			d = d.AddDate(1, 0, 0)
		}
		return d.Format(dateFormat), nil
	case strings.HasPrefix(repeat, "d "):
		var days int
		_, err := fmt.Sscanf(repeat, "d %d", &days)
		if err != nil || days <= 0 || days > 400 {
			return "", fmt.Errorf("invalid repeat rule: %v", repeat)
		}
		d = d.AddDate(0, 0, days)
		for d.Format(dateFormat) <= now.Format(dateFormat) {
			d = d.AddDate(0, 0, days)
		}
		return d.Format(dateFormat), nil
	case strings.HasPrefix(repeat, "w "):
		var daysOfWeek string
		_, err := fmt.Sscanf(repeat, "w %s", &daysOfWeek)
		if err != nil {
			return "", fmt.Errorf("invalid repeat rule: %v", repeat)
		}
		days := strings.Split(daysOfWeek, ",")
		for _, day := range days {
			dayInt, err := strconv.Atoi(day)
			if err != nil || dayInt < 1 || dayInt > 7 {
				return "", fmt.Errorf("invalid repeat rule: %v", repeat)
			}
		}
		d = d.AddDate(0, 0, 1)
		for d.Format(dateFormat) <= now.Format(dateFormat) {
			d = d.AddDate(0, 0, 1)
		}
		return d.Format(dateFormat), nil
	default:
		return "", fmt.Errorf("unsupported repeat rule: %v", repeat)
	}
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeatStr := r.FormValue("repeat")

	now, err := time.Parse(dateFormat, nowStr)
	if err != nil {
		http.Error(w, "invalid now date format", http.StatusBadRequest)
		return
	}

	nextDate, err := NextDate(now, dateStr, repeatStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(nextDate))
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Error decoding JSON: %v"}`, err), http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		http.Error(w, `{"error": "Title is required"}`, http.StatusBadRequest)
		return
	}

	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format(dateFormat)
	} else {
		parsedDate, err := time.Parse(dateFormat, task.Date)
		if err != nil {
			http.Error(w, `{"error": "Invalid date format"}`, http.StatusBadRequest)
			return
		}
		if parsedDate.Format(dateFormat) < now.Format(dateFormat) {
			if task.Repeat == "" {
				task.Date = now.Format(dateFormat)
			} else {
				nextDate, err := NextDate(now, task.Date, task.Repeat)
				if err != nil {
					http.Error(w, fmt.Sprintf(`{"error": "Invalid repeat rule: %v"}`, err), http.StatusBadRequest)
					return
				}
				task.Date = nextDate
			}
		}
	}

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Failed to insert task: %v"}`, err), http.StatusInternalServerError)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Failed to retrieve last insert ID: %v"}`, err), http.StatusInternalServerError)
		return
	}

	task.ID = strconv.FormatInt(id, 10)
	response := map[string]interface{}{"id": task.ID}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(response)
}

func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT 50")
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Failed to query tasks: %v"}`, err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tasks []map[string]string
	for rows.Next() {
		var task Task
		var id int64
		err := rows.Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "Failed to scan task: %v"}`, err), http.StatusInternalServerError)
			return
		}
		task.ID = strconv.FormatInt(id, 10)
		taskMap := map[string]string{
			"id":      task.ID,
			"date":    task.Date,
			"title":   task.Title,
			"comment": task.Comment,
			"repeat":  task.Repeat,
		}
		tasks = append(tasks, taskMap)
	}

	// Check for errors from iterating over rows.
	if err = rows.Err(); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Error iterating over rows: %v"}`, err), http.StatusInternalServerError)
		return
	}

	if tasks == nil {
		tasks = []map[string]string{}
	}

	response := map[string][]map[string]string{"tasks": tasks}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(response)
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, `{"error": "Не указан идентификатор"}`, http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error": "Неверный формат идентификатора"}`, http.StatusBadRequest)
		return
	}

	var task Task
	err = db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error": "Задача не найдена"}`, http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf(`{"error": "Ошибка при поиске задачи: %v"}`, err), http.StatusInternalServerError)
		}
		return
	}

	response := map[string]string{
		"id":      task.ID,
		"date":    task.Date,
		"title":   task.Title,
		"comment": task.Comment,
		"repeat":  task.Repeat,
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(response)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Error decoding JSON: %v"}`, err), http.StatusBadRequest)
		return
	}

	if task.ID == "" {
		http.Error(w, `{"error": "ID is required"}`, http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		http.Error(w, `{"error": "Title is required"}`, http.StatusBadRequest)
		return
	}

	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format(dateFormat)
	} else {
		parsedDate, err := time.Parse(dateFormat, task.Date)
		if err != nil {
			http.Error(w, `{"error": "Invalid date format"}`, http.StatusBadRequest)
			return
		}
		if parsedDate.Format(dateFormat) < now.Format(dateFormat) {
			if task.Repeat == "" {
				task.Date = now.Format(dateFormat)
			} else {
				nextDate, err := NextDate(now, task.Date, task.Repeat)
				if err != nil {
					http.Error(w, fmt.Sprintf(`{"error": "Invalid repeat rule: %v"}`, err), http.StatusBadRequest)
					return
				}
				task.Date = nextDate
			}
		}
	}

	id, err := strconv.ParseInt(task.ID, 10, 64)
	if err != nil {
		http.Error(w, `{"error": "Invalid ID format"}`, http.StatusBadRequest)
		return
	}

	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, id)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Failed to update task: %v"}`, err), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Failed to retrieve rows affected: %v"}`, err), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, `{"error": "Задача не найдена"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(map[string]string{})
}

func taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Method Not Allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, `{"error": "Не указан идентификатор"}`, http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error": "Неверный формат идентификатора"}`, http.StatusBadRequest)
		return
	}

	var task Task
	err = db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error": "Задача не найдена"}`, http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf(`{"error": "Ошибка при поиске задачи: %v"}`, err), http.StatusInternalServerError)
		}
		return
	}

	if task.Repeat == "" {
		_, err = db.Exec("DELETE FROM scheduler WHERE id = ?", id)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "Failed to delete task: %v"}`, err), http.StatusInternalServerError)
			return
		}
	} else {
		now := time.Now()
		nextDate, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "Failed to calculate next date: %v"}`, err), http.StatusInternalServerError)
			return
		}

		_, err = db.Exec("UPDATE scheduler SET date = ? WHERE id = ?", nextDate, id)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "Failed to update task: %v"}`, err), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(map[string]string{})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, `{"error": "Не указан идентификатор"}`, http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error": "Неверный формат идентификатора"}`, http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM scheduler WHERE id = ?", id)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Failed to delete task: %v"}`, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(map[string]string{})
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		http.Error(w, `{"error": "Method Not Allowed"}`, http.StatusMethodNotAllowed)
	}
}

func main() {
	if err := initDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	http.Handle("/", http.FileServer(http.Dir("./web")))
	http.HandleFunc("/api/nextdate", nextDateHandler)
	http.HandleFunc("/api/tasks", getTasksHandler)
	http.HandleFunc("/api/task", taskHandler)          // этот хэндлер решит, куда отправить — на создание задачи, получение по айди или редактирование
	http.HandleFunc("/api/task/done", taskDoneHandler) // Маршрут для выполнения задачи

	port := 7540
	addr := ":" + strconv.Itoa(port)
	log.Printf("Starting server on port %d...\n", port)
	log.Fatal(http.ListenAndServe(addr, nil))
}
