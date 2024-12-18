package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
)

type Task struct {
	ID       int64  `json:"id"`
	Date     string `json:"date"`
	Title    string `json:"title"`
	Comment  string `json:"comment"`
	Repeat   string `json:"repeat"`
	Done     bool   `json:"done"`
	NextDate string `json:"-"`
}

type Response struct {
	ID       int64  `json:"id,omitempty"`
	NextDate string `json:"next_date,omitempty"`
	Tasks    []Task `json:"tasks,omitempty"`
	Error    string `json:"error,omitempty"`
}

type Config struct {
	WebDir string
	Port   string
	DBPath string
}

const (
	webDir   = "./web"
	dbName   = "scheduler.db"
	tableSQL = `CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date TEXT,
		title TEXT,
		comment TEXT,
		repeat TEXT,
		done INTEGER DEFAULT 0
	);`
)

var config Config

func main() {
	// закрепляем путь к БД, фронту и порт
	config = Config{
		WebDir: "./web",
		Port:   ":7540",
		DBPath: "./scheduler.db",
	}

	// Подлкючаем БД или создаём её
	appPath, err := os.Executable()
	if err != nil {
		log.Fatalf("Не удалось получить путь приложения: %v\n", err)
	}
	dbFile := filepath.Join(filepath.Dir(appPath), dbName)
	db, err := openDB(dbFile)
	if err != nil {
		log.Fatalf("Ошибка открытия базы данных: %v\n", err)
	}
	defer db.Close()

	if !tableExists(db, "scheduler") {
		createTable(db)
	}

	// Настройка маршрутов и запуск сервера
	router := setupRouter()
	log.Printf("Запуск сервера на порту %s...\n", config.Port)
	log.Fatal(http.ListenAndServe(config.Port, router))
}

func setupRouter() *chi.Mux {
	router := chi.NewRouter()

	fs := http.FileServer(http.Dir(config.WebDir))
	router.Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	})

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(config.WebDir, "index.html"))
	})

	// API маршруты
	router.Get("/api/nextdate", nextDateHandler)
	router.Route("/api/task", func(r chi.Router) {
		r.With(PostRequestValidation).Post("/", withDB(createTaskHandler))
		r.With(GetRequestValidation).Get("/{id}", withDB(getTaskHandler))
		r.With(DeleteRequestValidation).Delete("/{id}", withDB(deleteTaskHandler))
		r.With(PutRequestValidation).Put("/{id}", withDB(updateTaskHandler))
	})
	router.Get("/api/tasks", withDB(getTasksHandler))
	router.Post("/api/task/done", withDB(markTaskDoneHandler))

	return router
}

// ------------------------ База данных ------------------------

func openDB(file string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть базу данных: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("не удалось подключиться к базе данных: %w", err)
	}
	return db, nil
}

func tableExists(db *sql.DB, tableName string) bool {
	row := db.QueryRow(fmt.Sprintf(`SELECT name FROM sqlite_master WHERE type='table' AND name='%s';`, tableName))
	var name string
	err := row.Scan(&name)
	return err == nil && strings.EqualFold(name, tableName)
}

func createTable(db *sql.DB) {
	_, err := db.Exec(tableSQL)
	if err != nil {
		log.Fatalf("Ошибка создания таблицы: %v\n", err)
	}
	log.Println("Таблица 'scheduler' создана.")
}

func withDB(handler func(db *sql.DB, w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db, err := openDB(config.DBPath)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Database error")
			return
		}
		defer db.Close()
		handler(db, w, r)
	}
}

// ------------------------ Вспомогательные функции ------------------------

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, errMsg string) {
	writeJSON(w, status, map[string]string{"error": errMsg})
}

// ------------------------ Обработчики ------------------------

// создание задачи
func createTaskHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	query := `INSERT INTO scheduler (date, title, comment, repeat, done) VALUES (?, ?, ?, ?, 0)`
	result, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Ошибка при создании задачи")
		return
	}

	taskID, _ := result.LastInsertId()
	writeJSON(w, http.StatusOK, Response{ID: taskID})
}

// получение задачи
func getTaskHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")

	query := `SELECT id, date, title, comment, repeat, done FROM scheduler WHERE id = ?`
	row := db.QueryRow(query, taskID)

	var task Task
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat, &task.Done)
	if err != nil {
		writeError(w, http.StatusNotFound, "Задача не найдена")
		return
	}

	writeJSON(w, http.StatusOK, task)
}

// обновление задачи
func updateTaskHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, http.StatusBadRequest, "Неверный ввод")
		return
	}

	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ?, done = ? WHERE id = ?`
	_, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.Done, taskID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось обновить задачу")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

// удаление задачи
func deleteTaskHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")

	query := `DELETE FROM scheduler WHERE id = ?`
	_, err := db.Exec(query, taskID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось удалить задачу")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

// постановка статуса
func markTaskDoneHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var payload struct {
		ID   int64 `json:"id"`
		Done bool  `json:"done"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, "Неверный ввод")
		return
	}

	query := `UPDATE scheduler SET done = ? WHERE id = ?`
	_, err := db.Exec(query, payload.Done, payload.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось поставить статус задаче")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

// запрос на вычисление следующей даты
func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	date := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	nextDate, err := NextDate(now, date, repeat)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, Response{NextDate: nextDate})
}

// вычисление следующей даты
func NextDate(now time.Time, date, repeat string) (string, error) {
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return "", errors.New("неверный формат даты")
	}

	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", errors.New("неверный формат повторения")
	}

	switch parts[0] {
	case "d":
		interval, err := strconv.Atoi(parts[1])
		if err != nil || interval < 1 || interval > 400 {
			return "", errors.New("неверный интервал повторения, нужно выбрать от 1 до 400 дней")
		}
		for parsedDate.Before(now) || parsedDate.Equal(now) {
			parsedDate = parsedDate.AddDate(0, 0, interval)
		}
	case "y":
		for parsedDate.Before(now) || parsedDate.Equal(now) {
			parsedDate = parsedDate.AddDate(1, 0, 0)
		}
	default:
		return "", errors.New("неподдерживаемый тип повторения")
	}

	return parsedDate.Format("2006-01-02"), nil
}

// получение всех задач
func getTasksHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Извлечение всех задач из БД
	query := `SELECT id, date, title, comment, repeat, done FROM scheduler`
	rows, err := db.Query(query)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Ошибка чтения задач")
		return
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat, &task.Done); err != nil {
			writeError(w, http.StatusInternalServerError, "Ошибка чтения задач")
			return
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		writeError(w, http.StatusInternalServerError, "Ошибка чтения задач")
		return
	}

	writeJSON(w, http.StatusOK, Response{Tasks: tasks})
}

// PostRequestValidation проверяет, что в запросе присутствуют обязательные поля
func PostRequestValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var task Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			writeError(w, http.StatusBadRequest, "Неверный ввод")
			return
		}

		if task.Date == "" || task.Title == "" {
			writeError(w, http.StatusBadRequest, "Дата и название задачи обязательны")
			return
		}

		// Если валидация пройдена, передаём управление следующему обработчику
		next.ServeHTTP(w, r)
	})
}

// GetRequestValidation проверяет, что ID задачи присутствует в URL
func GetRequestValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if _, err := strconv.Atoi(id); err != nil {
			writeError(w, http.StatusBadRequest, "Неверный ID задачи")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// DeleteRequestValidation проверяет, что ID задачи корректен для удаления
func DeleteRequestValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if _, err := strconv.Atoi(id); err != nil {
			writeError(w, http.StatusBadRequest, "Неверный ID задачи")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// PutRequestValidation проверяет, что тело запроса корректно для обновления задачи
func PutRequestValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var task Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			writeError(w, http.StatusBadRequest, "Неверный ввод")
			return
		}

		if task.Date == "" || task.Title == "" {
			writeError(w, http.StatusBadRequest, "Дата и название задачи обязательны")
			return
		}

		next.ServeHTTP(w, r)
	})
}
