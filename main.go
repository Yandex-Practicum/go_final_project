package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}
	dbPath := getDbPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err := initDb(db); err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}

	webDir := "./web"
	fileServer := http.FileServer(http.Dir(webDir))
	http.Handle("/", fileServer)

	// API рабы с аутентификацией
	http.HandleFunc("/api/nextdate", nextDateHandler)
	http.HandleFunc("/api/task", authMiddleware(makeHandler(taskHandler, db)))
	http.HandleFunc("/api/tasks", authMiddleware(makeHandler(tasksHandler, db)))
	http.HandleFunc("/api/task/done", authMiddleware(makeHandler(taskDoneHandler, db)))
	http.HandleFunc("/api/signin", makeHandler(signInHandler, db))

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

// Обработчик для /api/task
func taskHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r, db)
	case http.MethodGet:
		getTaskHandler(w, r, db)
	case http.MethodPut:
		editTaskHandler(w, r, db)
	case http.MethodDelete:
		deleteTaskHandler(w, r, db)
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

func getDbPath() string {
	dbPath := os.Getenv("TODO_DBFILE")
	if dbPath == "" {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		dbPath = filepath.Join(cwd, "scheduler.db")
	}
	return dbPath
}
func initDb(db *sql.DB) error {
	_, err := os.Stat(getDbPath())
	install := false
	if os.IsNotExist(err) {
		install = true
	}
	if install {
		createTableQuery := `
		CREATE TABLE scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date TEXT NOT NULL,
			title TEXT NOT NULL,
			comment TEXT,
			repeat TEXT(128)
		);
		CREATE INDEX idx_date ON scheduler(date);
		`
		_, err := db.Exec(createTableQuery)
		if err != nil {
			return err
		}
	}
	return nil
}
