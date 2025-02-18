package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// DBFile - путь к файлу базы данных.
var DBFile = "./scheduler.db"

func initDB() error {
	// Получаем путь к исполняемому файлу
	appPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %v", err)
	}

	// Определяем полный путь к файлу базы данных
	dbFile := filepath.Join(filepath.Dir(appPath), DBFile)

	// Проверяем, существует ли файл базы данных
	_, err = os.Stat(dbFile)
	var install bool
	if err != nil {
		if os.IsNotExist(err) {
			install = true
		} else {
			return fmt.Errorf("failed to check database file: %v", err)
		}
	}

	// Открываем или создаем базу данных
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	if install {
		// Если база данных не существует, создаем таблицу и индекс
		_, err = db.Exec(`
            CREATE TABLE scheduler (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                date TEXT NOT NULL,
                title TEXT NOT NULL,
                comment TEXT,
                repeat TEXT CHECK (LENGTH(repeat) <= 128)
            );
            CREATE INDEX idx_date ON scheduler(date);
        `)
		if err != nil {
			return fmt.Errorf("failed to create table: %v", err)
		}
		log.Println("Database and table created successfully.")
	} else {
		log.Println("Database already exists.")
	}

	return nil
}
