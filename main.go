package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// createDB инициализирует базу данных SQLite и создает необходимые таблицы.
func createDB() {
	db, err := sql.Open("sqlite3", "./scheduler.db")
	if err != nil {
		log.Fatalf("Ошибка при открытии базы данных: %v", err)
	}
	defer db.Close()

	commands := []string{
		`CREATE TABLE IF NOT EXISTS scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date CHAR(8),
			title TEXT,
			comment TEXT,
			repeat VARCHAR(128)
		)`,
		"CREATE INDEX IF NOT EXISTS indexdate ON scheduler (date)",
	}

	for _, cmd := range commands {
		if _, err := db.Exec(cmd); err != nil {
			log.Fatalf("Ошибка при выполнении команды: %s, ошибка: %v", cmd, err)
		}
	}
}

// taskHandler обрабатывает HTTP-запросы для операций, связанных с задачами.
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTask(w, r)
	case http.MethodGet:
		if id := r.URL.Query().Get("id"); id != "" {
			getTaskByID(w, r, id)
		} else {
			getTasks(w, r)
		}
	default:
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
	}
}

// addTask добавляет новую задачу в базу данных.
func addTask(w http.ResponseWriter, r *http.Request) {
	var task struct {
		Date    string `json:"date"`
		Title   string `json:"title"`
		Comment string `json:"comment"`
		Repeat  string `json:"repeat"`
	}

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Неверный запрос: "+err.Error(), http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		http.Error(w, "Требуется указать название", http.StatusBadRequest)
		return
	}

	if task.Date == "" {
		task.Date = time.Now().Format("20060102")
	} else if _, err := time.Parse("20060102", task.Date); err != nil {
		http.Error(w, "Неверный формат даты, ожидается YYYYMMDD", http.StatusBadRequest)
		return
	}

	if task.Repeat != "" {
		if _, err := NextDate(time.Now(), task.Date, task.Repeat); err != nil {
			http.Error(w, "Некорректное правило повторения: "+err.Error(), http.StatusBadRequest)
			return
		}
	}

	db, err := sql.Open("sqlite3", "./scheduler.db")
	if err != nil {
		http.Error(w, "Ошибка сервера: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO scheduler(date, title, comment, repeat) VALUES (?, ?, ?, ?)")
	if err != nil {
		http.Error(w, "Ошибка сервера: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		http.Error(w, "Ошибка при вставке задачи: "+err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		http.Error(w, "Ошибка при получении id: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(struct {
		ID int64 `json:"id"`
	}{ID: id})
}

// getTasks извлекает задачи из базы данных и возвращает их в формате JSON.
func getTasks(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./scheduler.db")
	if err != nil {
		http.Error(w, "Ошибка сервера: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC")
	if err != nil {
		http.Error(w, "Ошибка при получении задач: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tasks []map[string]string
	for rows.Next() {
		var id, date, title, comment, repeat string
		if err := rows.Scan(&id, &date, &title, &comment, &repeat); err != nil {
			http.Error(w, "Ошибка при сканировании задачи: "+err.Error(), http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, map[string]string{
			"id":      id,
			"date":    date,
			"title":   title,
			"comment": comment,
			"repeat":  repeat,
		})
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Ошибка при чтении данных задач: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(map[string]interface{}{"tasks": tasks})
}

// getTaskByID извлекает задачу по идентификатору из базы данных и возвращает ее в формате JSON.
func getTaskByID(w http.ResponseWriter, r *http.Request, id string) {
	db, err := sql.Open("sqlite3", "./scheduler.db")
	if err != nil {
		http.Error(w, "Ошибка сервера: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var task struct {
		ID      string `json:"id"`
		Date    string `json:"date"`
		Title   string `json:"title"`
		Comment string `json:"comment"`
		Repeat  string `json:"repeat"`
	}

	err = db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			json.NewEncoder(w).Encode(map[string]string{"error": "Задача не найдена"})
		} else {
			http.Error(w, "Ошибка при получении задачи: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(task)
}

// markTaskDone обрабатывает POST-запрос для выполнения задачи.
func markTaskDone(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(map[string]string{"error": "Не указан идентификатор"})
		return
	}

	db, err := sql.Open("sqlite3", "./scheduler.db")
	if err != nil {
		http.Error(w, "Ошибка сервера: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var task struct {
		ID      string `json:"id"`
		Date    string `json:"date"`
		Title   string `json:"title"`
		Comment string `json:"comment"`
		Repeat  string `json:"repeat"`
	}

	err = db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			json.NewEncoder(w).Encode(map[string]string{"error": "Задача не найдена"})
		} else {
			http.Error(w, "Ошибка при получении задачи: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if task.Repeat == "" {
		_, err = db.Exec("DELETE FROM scheduler WHERE id = ?", id)
		if err != nil {
			http.Error(w, "Ошибка при удалении задачи: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(map[string]string{})
		return
	}

	nextDate, err := NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		http.Error(w, "Ошибка при вычислении следующей даты: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("UPDATE scheduler SET date = ? WHERE id = ?", nextDate, id)
	if err != nil {
		http.Error(w, "Ошибка при обновлении даты задачи: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(map[string]string{})
}

func main() {
	if _, err := os.Stat("./scheduler.db"); os.IsNotExist(err) {
		createDB()
	}

	http.HandleFunc("/api/task", taskHandler)
	http.HandleFunc("/api/task/done", markTaskDone)
	http.Handle("/", http.FileServer(http.Dir("./web")))
	log.Fatal(http.ListenAndServe(":7540", nil))
}
