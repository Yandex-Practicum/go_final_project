package dbase

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(256) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat VARCHAR(128) NOT NULL DEFAULT ""

);
CREATE INDEX date_index ON scheduler (date)
`

var DB *sql.DB

func Init(dbFile string) error {
	var err error
	_, err = os.Stat(dbFile)
	flag := os.IsNotExist(err)

	// Открываем соединение с БД
	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("failed to open DB: %v", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("error ping DB: %v", err)
	}
	if flag {
		if _, err = DB.Exec(schema); err != nil {
			return fmt.Errorf("error create schema: %v", err)
		}
	}
	return nil
}
