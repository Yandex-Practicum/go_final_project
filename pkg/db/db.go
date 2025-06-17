package db

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

var (
	database *sqlx.DB
	DBFile   = "scheduler.db" // Дефолтный путь к файлу базы данных
)

func Init(dbFile string) error {
	var err error
	database, err = sqlx.Connect("sqlite", dbFile)
	if err != nil {
		return err
	}

	// Создание таблицы, если она не существует
	schema := `
	CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date TEXT,
		title TEXT,
		comment TEXT,
		repeat TEXT
	);
	`
	_, err = database.Exec(schema)
	return err
}

func Close() {
	if database != nil {
		err := database.Close()
		if err != nil {
			log.Printf("Ошибка при закрытии базы данных: %v", err)
		}
	}
}
