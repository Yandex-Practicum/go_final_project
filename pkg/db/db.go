package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

var db *sql.DB

// SQL схема для создания таблицы и индекса
const schema = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(255) NOT NULL DEFAULT "",
    comment TEXT,
    repeat VARCHAR(128)
);

CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);
`

// Init открывает базу данных и при необходимости создаёт таблицу
func Init(dbFile string) error {
	_, err := os.Stat(dbFile)
	install := false
	if err != nil {
		if os.IsNotExist(err) {
			install = true
		} else {
			return fmt.Errorf("не удалось проверить файл базы данных: %w", err)
		}
	}

	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("не удалось открыть базу данных: %w", err)
	}

	if install {
		if _, err := db.Exec(schema); err != nil {
			return fmt.Errorf("не удалось создать таблицу: %w", err)
		}
	}

	return nil
}

// Close закрывает базу данных
func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

// DB возвращает глобальный объект базы данных
func DB() *sql.DB {
	return db
}
