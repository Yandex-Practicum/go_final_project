package config

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

var db *sql.DB

func GetDB() *sql.DB {
	return db
}

func InitializeDatabase() error {
	// Получаем значение переменной окружения TODO_DBFILE
	dbFile := os.Getenv("TODO_DBFILE")

	// Проверяем, существует ли файл базы данных
	_, err := os.Stat(dbFile)
	install := false
	if err != nil {
		if os.IsNotExist(err) {
			install = true
		} else {
			return fmt.Errorf("Файл базы данных не существует: %v", err)
		}
	}

	// Подключаемся к базе данных
	db, err = sql.Open("sqlite3", dbFile)
	if err != nil {
		return fmt.Errorf("Ошибка подключения: %v", err)
	}

	// Если база данных новая, создаём таблицу
	if install {
		log.Println("Создаётся новая база данных...")
		createTableQuery := `
		CREATE TABLE scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date TEXT NOT NULL,
			title TEXT NOT NULL,
			comment TEXT,
			repeat TEXT CHECK(length(repeat) <= 128)
		);
		CREATE INDEX idx_date ON scheduler(date);
		`
		_, err := db.Exec(createTableQuery)
		if err != nil {
			return fmt.Errorf("Ошибка при создании таблицы: %v", err)
		}
		log.Println("Таблица scheduler создана.")
	}

	log.Println("Database connection established")
	return nil
}

func CloseDB() {
	if db != nil {
		err := db.Close()
		if err != nil {
			log.Printf("Error closing io connection: %v", err)
		} else {
			log.Println("Database connection closed")
		}
	}
}
