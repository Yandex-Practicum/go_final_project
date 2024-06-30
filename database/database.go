package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

const dbFileName = "scheduler.db"

// InitializeDatabase проверяет существование файла базы данных и создает таблицу, если необходимо
func InitializeDatabase() (*sql.DB, error) {
	if _, err := os.Stat(dbFileName); os.IsNotExist(err) {
		file, err := os.Create(dbFileName)
		if err != nil {
			return nil, err
		}
		file.Close()
	}

	db, err := sql.Open("sqlite3", dbFileName)
	if err != nil {
		return nil, err
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS scheduler (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
		date CHAR(8) NOT NULL DEFAULT "",
		title VARCHAR(128) NOT NULL DEFAULT "",
		comment TEXT NOT NULL DEFAULT "",
		repeat VARCHAR(128) NOT NULL DEFAULT ""
    );`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}

	// Создаем индекс по полю date для сортировки задач по дате
	createIndexSQL := `CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);`
	_, err = db.Exec(createIndexSQL)
	if err != nil {
		return nil, err
	}

	return db, nil
}
