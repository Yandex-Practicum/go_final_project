package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"my_education/go/go_final_project/internal/handlers"
	"net/http"
	"os"
)

const testPort = "7540"
const webDir = "./web"

// initializeDB создает таблицу scheduler и индекс по полю date, если их нет
func initializeDB(db *sql.DB) error {
	createTableQuery := `
 CREATE TABLE IF NOT EXISTS scheduler (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  date TEXT NOT NULL,
  title TEXT NOT NULL,
  comment TEXT,
  repeat TEXT(128)
 );
 CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);
 `
	_, err := db.Exec(createTableQuery)
	return err
}

// setupDatabase открывает подключение к базе данных и создает таблицу при необходимости
func setupDatabase(dbFile string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	var tableName string
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='scheduler'").Scan(&tableName)
	if err == sql.ErrNoRows {
		log.Println("Таблица 'scheduler' не найдена, создаем таблицу...")
		err = initializeDB(db)
		if err != nil {
			return nil, err
		}
		log.Println("Таблица 'scheduler' успешно создана.")
	} else if err != nil {
		return nil, err
	} else {
		log.Println("Таблица 'scheduler' уже существует.")
	}

	return db, nil
}

// startServer настраивает и запускает HTTP-сервер
func startServer(port string, db *sql.DB) error {
	fs := http.FileServer(http.Dir(webDir))
	http.Handle("/", fs)

	// Регистрация обработчиков для API
	http.HandleFunc("/api/nextdate", handlers.NextDateHandler)
	http.HandleFunc("/api/task", handlers.TaskHandler(db))
	http.Handle("/api/tasks", handlers.GetTasksHandler(db))
	http.HandleFunc("/api/task/done", handlers.MarkTaskDoneHandler(db))

	log.Printf("Сервер запущен - порт %s\n", port)
	return http.ListenAndServe(":"+port, nil)
}

func main() {
	dbFile := "scheduler.db"

	// Инициализия бд
	db, err := setupDatabase(dbFile)
	if err != nil {
		log.Fatalf("Ошибка при настройке базы данных: %v\n", err)
	}
	defer db.Close()

	// Установка порта для сервера
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = testPort
	}

	// Запускаем сервер
	if err := startServer(port, db); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v\n", err)
	}
}
