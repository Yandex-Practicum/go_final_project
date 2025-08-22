package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func CloseDB() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(255) NOT NULL,
    comment TEXT,
    repeat VARCHAR(128)
);

CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);
`

func Init(dbFile string) error {
	// Проверяем существование файла
	_, err := os.Stat(dbFile)
	var install bool
	if err != nil {
		install = true
	}

	// Открываем или создаем базу данных
	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("не удалось открыть базу данных: %w", err)
	}

	// Проверяем соединение
	err = db.Ping()
	if err != nil {
		return fmt.Errorf("не удалось установить соединение: %w", err)
	}

	// Если файл не существовал, создаем схему
	if install {
		_, err = db.Exec(schema)
		if err != nil {
			return fmt.Errorf("ошибка создания схемы: %w", err)
		}
		log.Println("База данных успешно создана")
	}

	return nil
}
