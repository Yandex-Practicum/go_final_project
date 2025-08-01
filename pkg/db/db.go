package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE scheduler  (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(256) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat VARCHAR(128) NOT NULL DEFAULT ""
);
CREATE INDEX scheduler_date_index ON scheduler (date);
`

var db *sql.DB

func Close() {
	db.Close()
}

func Init(dbFile string) error {
	// Проверяем существование файла БД
	_, err := os.Stat(dbFile)
	install := os.IsNotExist(err)

	db, err = sql.Open("sqlite", dbFile)

	if err != nil {
		return fmt.Errorf("can't open database: %w", err)
	}

	if install {
		if _, err = db.Exec(schema); err != nil {
			return fmt.Errorf("can't create schema: %w", err)
		}

	}
	return nil
}
