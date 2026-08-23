package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

// SQL запрос для создания основной таблицы планировщика.
const createTable = `
CREATE TABLE scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date CHAR(8) NOT NULL DEFAULT '',
	title VARCHAR(128) NOT NULL DEFAULT '',
	comment TEXT NOT NULL DEFAULT '',
	repeat VARCHAR(128) NOT NULL DEFAULT ''
); 
`

// Индекс по дате ускоряет выборку задач, отсортированных по дате.
const createIndex = `
CREATE INDEX idx_date ON scheduler(date);
`

var db *sql.DB

// Init открывает базу данных и при необходимости создает таблицу scheduler
// и индекс по дате.
// dbFile содержит путь к файлу базы данных
func Init(dbFile string) error {
	_, err := os.Stat(dbFile)

	var install bool

	if err != nil {
		if os.IsNotExist(err) {
		install = true
	} else {
		return fmt.Errorf("check database file: %w", err)
	}
	}

	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}

	if install == true {
		_, err = db.Exec(createTable)
		if err != nil {
			return fmt.Errorf("create table: %w", err)
		}
		_, err = db.Exec(createIndex)
		if err != nil {
			return fmt.Errorf("create index: %w", err)
		}
	}

	return nil
}