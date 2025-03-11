package db

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// InitDB открывает или создаёт базу данных и инициализирует таблицу scheduler.
func InitDB() (*sql.DB, error) {
	// Если переменная окружения TODO_DBFILE указана, используем её; иначе scheduler.db в папке приложения.
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		exe, err := os.Executable()
		if err != nil {
			return nil, err
		}
		dbFile = filepath.Join(filepath.Dir(exe), "scheduler.db")
	}

	// Если файл базы данных отсутствует, нужно создать таблицу.
	_, err := os.Stat(dbFile)
	install := err != nil

	database, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	if install {
		createTable := `
		CREATE TABLE scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date TEXT NOT NULL,
			title TEXT NOT NULL,
			comment TEXT,
			repeat TEXT
		);
		CREATE INDEX idx_date ON scheduler(date);
		`
		if _, err := database.Exec(createTable); err != nil {
			return nil, err
		}
		log.Println("База данных создана:", dbFile)
	}
	return database, nil
}
