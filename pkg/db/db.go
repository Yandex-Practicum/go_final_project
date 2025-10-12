package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

const DateFormat = "20060102"

var DB *sql.DB

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date CHAR(8) NOT NULL DEFAULT "",
	title VARCHAR(255) NOT NULL,
	comment TEXT,
	repeat VARCHAR(255)
	);
	
CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);
	`

// Init инициализирует базу данных
func Init(dbFile string) error {
	// Проверяем существование файла БД
	_, err := os.Stat(dbFile)
	install := os.IsNotExist(err)

	// Открываем базу данных
	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("ошибка открытия БД: %v", err)
	}

	// Проверяем соединение
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("ошибка подключения к БД: %v", err)
	}

	// Если файл не существовал, создаем
	if install {
		if _, err := DB.Exec(schema); err != nil {
			return fmt.Errorf("ошибка создания схемы: %v", err)
		}
		fmt.Printf("База данных создана: %s\n", dbFile)
	} else {
		fmt.Printf("База данных подключена: %s\n", dbFile)
	}

	return nil

}

// Close закрывает соединение с БД
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
