package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

// Соединение с Базой Данных
var db *sql.DB

// создаем таблицу scheduler
const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(255) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat VARCHAR(128) NOT NULL DEFAULT ""
);

CREATE INDEX IF NOT EXISTS scheduler_date ON scheduler(date);
`

// Init инициализация БД
func Init(dbFile string) error {

	//Проверка наличия файла
	_, err := os.Stat(dbFile)
	install := errors.Is(err, os.ErrNotExist)

	//Соединение с БД
	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("error opening sql %v", err)
	}

	//Создаем таблицу если ее не существует
	_, err = db.Exec(schema)
	if err != nil {
		return fmt.Errorf("error executing SQL query %v", err)
	}

	//Сообщение о создании таблицы
	if install {
		fmt.Println("A new database has been created:", dbFile)
	}

	//Удаляем из таблицы задачи без даты
	_, err = db.Exec(`DELETE FROM scheduler WHERE date = ''`)
	if err != nil {
		return fmt.Errorf("database cleanup error: %v", err)
	}

	return nil
}

// GetDB Доступ для других пакетов
func GetDB() *sql.DB {
	return db
}

// Close Для закрытия соединения
func Close() error {
	return db.Close()
}
