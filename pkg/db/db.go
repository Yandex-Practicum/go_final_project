package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

var (
	db     *sql.DB
	schema = `
        CREATE TABLE scheduler (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            date CHAR(8) NOT NULL DEFAULT "",
            title VARCHAR(255) NOT NULL,
            comment TEXT,
            repeat VARCHAR(128),
            UNIQUE(date, title) -- Уникальное сочетание даты и заголовка
        );
        CREATE INDEX idx_date ON scheduler(date);
    `
)

// Init инициализирует базу данных и создает таблицу, если она не существует
func Init(dbFile string) error {
	// Проверяем существование файла базы данных
	_, err := os.Stat(dbFile)
	install := false
	if err != nil {
		install = true
	}

	// Открываем базу данных
	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	// Если база данных новая, создаем таблицу и индекс
	if install {
		_, err = db.Exec(schema)
		if err != nil {
			return err
		}
	}

	return nil
}

// Close закрывает соединение с базой данных
func Close() error {
	return db.Close()
}
