package database

import (
	"database/sql"
	"log"
	"os"
)

func InitDB(dbFile string) (*sql.DB, error) {
	var err error
	// Определяем путь к базе данных
	if envDb := os.Getenv("TODO_DBFILE"); envDb != "" {
		dbFile = envDb
	}
	// Открываем базу данных
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	// Создаем таблицу и индекс, если файл базы данных отсутствует
	createTableQuery := `
        CREATE TABLE IF NOT EXISTS scheduler (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            date TEXT NOT NULL,
            title TEXT NOT NULL,
            comment TEXT,
            repeat TEXT CHECK(LENGTH(repeat) <= 128)
        );
        CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);
        `
	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, err
	}
	log.Println("Database открыта.")
	return db, nil
}
