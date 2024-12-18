package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

const createTableSQL = `
CREATE TABLE IF NOT EXISTS scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date TEXT NOT NULL,
	title TEXT NOT NULL,
	comment TEXT,
	repeat TEXT
);`

const createIndexSQL = `
CREATE INDEX IF NOT EXISTS idx_date ON scheduler (date);`

func InitializeDatabase() (*sql.DB, error) {
	appPath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("ошибка определения пути приложения: %w", err)
	}
	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")

	install := false
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		install = true
	}

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия базы данных: %w", err)
	}

	if install {
		log.Println("Создаю новую базу данных...")

		_, err = db.Exec(createTableSQL)
		if err != nil {
			return nil, fmt.Errorf("ошибка создания таблицы: %w", err)
		}

		_, err = db.Exec(createIndexSQL)
		if err != nil {
			return nil, fmt.Errorf("ошибка создания индекса: %w", err)
		}

		log.Println("База данных успешно создана.")
	}

	return db, nil
}
