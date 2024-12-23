package config

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

func InitializeDatabase() *sql.DB {
	// Получаем значение переменной окружения TODO_DBFILE
	dbFile := os.Getenv("TODO_DBFILE")

	// Проверяем, существует ли файл базы данных
	_, err := os.Stat(dbFile)
	install := false
	if err != nil {
		if os.IsNotExist(err) {
			install = true
		} else {
			log.Fatal(err)
		}
	}

	// Подключаемся к базе данных
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal(err)
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
			log.Fatal("Ошибка при создании таблицы: ", err)
		}
		log.Println("Таблица scheduler создана.")
	}

	return db
}
