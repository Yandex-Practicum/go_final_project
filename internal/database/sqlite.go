package database

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"path/filepath"
)

func Check() {
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")
	_, err = os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}
	if install == true {
		fmt.Println("База данных не найдена. Создаем...")
		CreateDB(dbFile)
	}
	// если install равен true, после открытия БД требуется выполнить
	// sql-запрос с CREATE TABLE и CREATE INDEX

	//TODO Реализуйте возможность определять путь к файлу базы данных через переменную окружения. Для этого сервер должен получать значение переменной окружения TODO_DBFILE и использовать его в качестве пути к базе данных, если это не пустая строка.
}

func CreateDB(path string) {
	file, err := os.Create(path)
	if err != nil {
		fmt.Println("Ошибка создания файла:", err)
		return
	}
	defer file.Close()

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal(err)
	}

	// Создаем базу данных и индекс
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date,
			title TEXT NOT NULL,
			comment TEXT,
			repeat TEXT CHECK(length(repeat) <= 128)
		);
		CREATE INDEX idx_date ON scheduler(date);
	`); err != nil {
		log.Fatal(err)
	}

	fmt.Println("База данных успешно создана и настроена.")
	// Закрываем соединение с базой данных
	db.Close()
}
