package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

var now = time.Now()

// createTable создает таблицу SQLite
func createTable() {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date INTEGER NOT NULL DEFAULT 20060102,
	title TEXT NOT NULL DEFAULT "",
	comment TEXT NOT NULL DEFAULT "",
	repeat TEXT CHECK (LENGTH(repeat) <= 128) NOT NULL DEFAULT "");`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = db.Exec("CREATE INDEX IF NOT EXISTS scheduler_date ON scheduler (date);")
	if err != nil {
		fmt.Println(err)
		return
	}
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("ошибка загрузки файла .env")
	}
}

func main() {
	// Подключение к БД
	dbPath := os.Getenv("TODO_DBFILE")
	if dbPath == "" {
		dbPath = "scheduler.db"
	}

	dbFile := filepath.Join(dbPath)
	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	db, err = sql.Open("sqlite3", "scheduler.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	if install {
		createTable()
	}

	r := chi.NewRouter()

	// Маршрут для статических файлов
	r.Handle("/*", http.StripPrefix("/", http.FileServer(http.Dir("web"))))

	// Маршруты API
	r.Get("/api/nextdate", NextDateHandler)
	r.Post("/api/task", Auth(PostTask))
	r.Get("/api/tasks", Auth(GetTasks))
	r.Get("/api/task", Auth(GetTasks))
	r.Put("/api/task", Auth(PutTask))
	r.Post("/api/task/done", Auth(TaskDone))
	r.Delete("/api/task", Auth(DeleteTask))
	r.Post("/api/signin", SignIn)

	// Настройка порта
	port := os.Getenv("TODO_PORT")

	if port == "" {
		port = "0.0.0.0:7540"
	}

	adress := "0.0.0.0:" + port

	// Запуск сервера
	if err := http.ListenAndServe(adress, r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
