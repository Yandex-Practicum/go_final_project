package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

// глобальная переменная = хранит подключение к бд
var db *sql.DB

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

// Init инициализирует базу данных и создает таблицу, если её нет
func InitDB(dbFile string) error {
	// Проверка, существует ли база данных
	_, err := os.Stat(dbFile)
	if err != nil {
		// Если база данных не существует, выводим сообщение
		fmt.Println("Таблица не найдена, создаем...")
	}

	// Открываем базу данных
	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("невозможно открыть таблицу: %v", err)
	}

	// Выполняем создание таблицы и индекса, если их нет
	_, err = db.Exec(schema)
	if err != nil {
		return fmt.Errorf("ошибка выполнения schema: %v", err)
	}

	return nil
}

// Функция для получения глобального подключения к базе данных
func GetDB() *sql.DB {
	return db
}
