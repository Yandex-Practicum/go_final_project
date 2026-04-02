package db

import (
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(255) NOT NULL DEFAULT "",
    comment TEXT,
    repeat VARCHAR(128) NOT NULL DEFAULT ""
);

CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);
`

func Init(dbFile string) error {

	_, err := os.Stat(dbFile)
	needInit := os.IsNotExist(err)

	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	if err = DB.Ping(); err != nil {
		return err
	}

	log.Println("Проверка/создание таблицы scheduler...")
	_, err = DB.Exec(schema)
	if err != nil {
		return err
	}

	if needInit {
		log.Printf("Создана новая база данных: %s", dbFile)
	} else {
		log.Printf("Использую существующую базу данных: %s", dbFile)
	}

	return nil
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
