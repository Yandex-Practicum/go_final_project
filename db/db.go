package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// InitDB открывает или создаёт базу данных по указанному пути,
// создаёт таблицу scheduler и индекс по полю date.
func InitDB(dbPath string) (*sql.DB, error) {
	log.Printf("Создание базы данных")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия БД: %w", err)
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date TEXT NOT NULL,
		title TEXT NOT NULL,
		comment TEXT,
		repeat TEXT
	);
	`
	_, err = db.Exec(createTable)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания таблицы: %w", err)
	}

	createIndex := `CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);`
	_, err = db.Exec(createIndex)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания индекса: %w", err)
	}

	return db, nil
}
