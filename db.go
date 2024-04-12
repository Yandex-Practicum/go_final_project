package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func getDBFilePath() string {
	dbFilePath := os.Getenv("TODO_DBFILE")
	return dbFilePath
}

func connectDB() (*sql.DB, error) {
	dbFilePath := getDBFilePath()
	if dbFilePath == "" {
		dbFilePath = "scheduler.db"
	}

	db, err := sql.Open("sqlite", dbFilePath)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func createTable(db *sql.DB) error {
	createStmt := `
		CREATE TABLE IF NOT EXISTS scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date TEXT,
			title TEXT,
			comment TEXT,
			repeat TEXT
		);
		CREATE INDEX IF NOT EXISTS idx_date ON scheduler (date);
	`
	fmt.Println("Создание таблицы перед тестированием...")
	_, err := db.Exec(createStmt)
	if err != nil {
		return err
	}
	fmt.Println("Таблица успешно создана перед тестированием.")
	return nil
}

func InitializeDB() {
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")

	newFunction1(err, dbFile)

	dbFilePath := os.Getenv("TODO_DBFILE")
	if dbFilePath != "" {
		fmt.Printf("Используется пользовательский путь к базе данных: %s\n", dbFilePath)
	}

	db, err := connectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = createTable(db)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("База данных и таблица успешно созданы.")
}

func newFunction1(err error, dbFile string) {
	_, err = os.Stat(dbFile)
}
