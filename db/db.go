package db

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// InitDB удаляет существующий файл базы данных (если он есть) и создаёт новую базу в указанном месте.
func InitDB() (*sql.DB, error) {
	// Если переменная окружения TODO_DBFILE указана, используем её; иначе база будет храниться в корневой директории проекта.
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	} else {
		absPath, err := filepath.Abs(dbFile)
		if err == nil {
			dbFile = absPath
		}
	}

	// Если файл базы данных существует, удаляем его.
	if _, err := os.Stat(dbFile); err == nil {
		if err := os.Remove(dbFile); err != nil {
			return nil, err
		}
	}

	database, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	// Создаем таблицу, если её нет.
	createTableStmt := `
	CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date TEXT NOT NULL,
		title TEXT NOT NULL,
		comment TEXT,
		repeat TEXT
	);`
	if _, err := database.Exec(createTableStmt); err != nil {
		return nil, err
	}

	// Создаем индекс, если его нет.
	createIndexStmt := `CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);`
	if _, err := database.Exec(createIndexStmt); err != nil {
		return nil, err
	}

	log.Println("База данных инициализирована:", dbFile)
	return database, nil
}
