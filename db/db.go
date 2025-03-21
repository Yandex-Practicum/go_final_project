package db

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite" // Импорт драйвера SQLite
)

// InitDB инициализирует базу данных и создает необходимые таблицы
func InitDB(isTest bool) (*sql.DB, error) {
	var dbFile string

	if isTest {
		dbFile = filepath.Join(os.TempDir(), "test_scheduler.db")
	} else {
		dbFile = "./scheduler.db"
	}

	if isTest {
		if _, err := os.Stat(dbFile); err == nil {
			os.Remove(dbFile)
		}
	}

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, err
	}

	var tableExists string
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='scheduler';").Scan(&tableExists)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// Создаем таблицу, если она не существует
	createTableSQL := `CREATE TABLE IF NOT EXISTS scheduler (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        date TEXT NOT NULL,
        title TEXT NOT NULL,
        comment TEXT,
        repeat TEXT
    );`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}

	// Создаем индекс
	createIndexSQL := `CREATE INDEX IF NOT EXISTS idx_date ON scheduler (date);`
	_, err = db.Exec(createIndexSQL)
	if err != nil {
		return nil, err
	}

	return db, nil
}
