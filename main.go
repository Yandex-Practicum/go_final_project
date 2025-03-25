package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type Response struct {
	ID    string `json:"id,omitempty"`
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

type TasksResponse struct {
	Tasks []Task `json:"tasks"`
}

type Password struct {
	Pass string `json:"password"`
}

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

func addTaskRules(t *Task) (int, error) {
	if t.Date == "" {
		t.Date = time.Now().Format(timeFormat)
	}
	_, err := time.Parse(timeFormat, t.Date)
	if err != nil {
		return 0, fmt.Errorf("дата представлена в формате, отличном от 20060102")
	}

	if t.Date < time.Now().Format(timeFormat) {
		if t.Repeat == "" {
			t.Date = time.Now().Format(timeFormat)
		} else {
			taskDay, err := DateParse(now, t.Date, t.Repeat)
			if err != nil {
				return 0, fmt.Errorf("ошибка при парсинге даты: %v", err)
			}
			t.Date = taskDay
		}
	}

	dateInt, err := strconv.Atoi(t.Date)
	if err != nil {
		return 0, fmt.Errorf("не удалось конвертировать дату в число")
	}

	if t.Title == "" {
		return 0, fmt.Errorf("не указан заголовок задачи")
	}

	return dateInt, nil
}

func main() {
	// appPath, err := os.Executable()
	// if err != nil {
	// 	log.Fatal(err)
	// }

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

	r.Handle("/*", http.StripPrefix("/", http.FileServer(http.Dir("web"))))
	r.Get("/api/nextdate", NextDateHandler)
	r.Post("/api/task", Auth(PostTask))
	r.Get("/api/tasks", Auth(GetTasks))
	r.Get("/api/task", Auth(GetTasks))
	r.Put("/api/task", Auth(PutTask))
	r.Post("/api/task/done", Auth(TaskDone))
	r.Delete("/api/task", Auth(DeleteTask))
	r.Post("/api/signin", SignIn)

	port := os.Getenv("TODO_PORT")

	if port == "" {
		port = "0.0.0.0:7540"
	}

	adress := "0.0.0.0:" + port

	if err := http.ListenAndServe(adress, r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
