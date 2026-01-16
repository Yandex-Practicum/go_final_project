package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

var db *sql.DB

const schema = `CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date CHAR(8) NOT NULL DEFAULT "",
		title VARCHAR NOT NULL,
		comment TEXT,
		repeat VARCHAR NOT NULL CHECK(length(repeat) <= 128)
	)`

// Init открывает базу данных или, при необходимости, создает ее
func Init(dbFile string) error {

	// Проверка на передачу пустой строки
	if dbFile == "" {
		return fmt.Errorf("Имя файла Базы данных не может быть пустым")
	}

	// Я не стал делать проверку на существование файла, так как sql.Open в любом случае
	// либо его создаст, либо откроет, а CREATE IF NOT EXISTS сделает то же самое с
	// таблицей

	var err error
	// Создаем или открываем БД
	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("Невозможно создать или открыть БД: %w", err)
	}

	// Создаем таблицу если ее нет
	_, err = db.Exec(schema)
	if err != nil {
		db.Close()
		db = nil
		return fmt.Errorf("Не удалось создать таблицу: %w", err)
	}

	return nil
}

func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}
