package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// Инициализация БД
func initDB() error {
	dbPath := getDBPath()
	
	// Для тестов: удаляем старую БД
	if os.Getenv("TEST_MODE") == "1" {
		os.Remove(dbPath)
	}

	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	// Создаём таблицы (IF NOT EXISTS для надёжности)
	if err := createTables(); err != nil {
		return err
	}

	log.Println("DB initialized at:", dbPath)
	return nil
}

func getDBPath() string {
	if path := os.Getenv("TODO_DBFILE"); path != "" {
		return path
	}
	return "scheduler.db" // Просто в текущей директории
}

func createTables() error {
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date TEXT NOT NULL,
			title TEXT NOT NULL,
			comment TEXT,
			repeat TEXT(128)
		);
		CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);
	`)
	return err
}