package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite" // Импортируем пакет SQLite
)

const bdName = "scheduler.db" // Имя файла базы данных
var db *sql.DB                // переменная БД

// createTable создает базу данных и таблицу scheduler
func createTable() error {
	// Открываем соединение с базой данных
	db, err := sql.Open("sqlite", bdName)
	if err != nil {
		return err // Возвращаем ошибку, если не удалось открыть базу данных
	}
	defer db.Close() // Закрываем соединение после завершения работы с базой данных

	// SQL-запрос для создания таблицы scheduler
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date TEXT NOT NULL,
		title TEXT NOT NULL,
		comment TEXT,
		repeat TEXT CHECK(length(repeat) <= 128)
	);
	CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date); -- Индекс по полю date
	`

	// Выполняем SQL-запрос для создания таблицы
	_, err = db.Exec(createTableSQL)
	return err
}

const defaultPort = "7540"
const webDir = "./web"

// Универсальный обработчик для GET и POST запросов к /api/task
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTaskHandler(w, r, db) // Получение задачи по ID
	case http.MethodPost:
		addTaskHandler(w, r, db) // Добавление задачи
	case http.MethodPut:
		updateTaskHandler(w, r, db) //Добавили обработку PUT-запроса
	case http.MethodDelete:
		deleteTaskHandler(w, r, db) // Удаление задачи
	default:
		JSONError(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func dbDB(db *sql.DB, handler func(http.ResponseWriter, *http.Request, *sql.DB)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, db)
	}
}

func main() {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = defaultPort
	}

	// проверка на наличие файла
	appPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	dbFile := filepath.Join(appPath, "scheduler.db")
	_, err = os.Stat(dbFile)
	if err != nil {
		fmt.Println("создаем")
		// Создаем базу данных и таблицу
		if err := createTable(); err != nil {
			log.Fatalf("Ошибка при создании базы данных: %v\n", err)
		}
	}

	// Подключаемся к БД
	db, err = sql.Open("sqlite", bdName) // Открываем БД один раз
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close() // Закрываем БД при завершении работы
	// Обработчик для API

	http.Handle("/", http.FileServer(http.Dir(webDir)))
	http.HandleFunc("/api/nextdate", NextDateHandler)
	http.HandleFunc("/api/task", taskHandler)

	http.HandleFunc("/api/tasks", dbDB(db, tasksHandler))
	http.HandleFunc("/api/task/done", dbDB(db, taskDoneHandler))

	// Запуск сервера
	fmt.Printf("Сервер запущен на [http://localhost:%s]\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Printf("Ошибка при запуске сервера: %v\n", err)
	}
}
