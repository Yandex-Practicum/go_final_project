package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const dbPath = "database/scheduler.db"

// СheckDB проверяет существет ли файл базы данных
// Создает, если не существует
func СheckDB() (*sql.DB, error) {
	appPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	//log.Printf("appPath=%v", appPath)
	dbFile := filepath.Join(appPath, dbPath)
	//log.Printf("dbFile=%v", dbFile)
	_, err = os.Stat(dbFile)
	//log.Printf("err=%v", err)
	var install bool
	if err != nil {
		install = true
	}
	// если install равен true, после создания БД требуется выполнить
	// sql-запрос с CREATE TABLE и CREATE INDEX
	if install {
		CreateDBFile(dbPath)
		database, err := sql.Open("sqlite", dbPath)
		if err != nil {
			log.Fatal(err)
		}
		defer database.Close()
		CreateTable(database)
	}
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// CreateDBFile создает файл базы данных в db/scheduler.db
func CreateDBFile(path string) {
	file, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()
	fmt.Printf("Table created successfully in %s!\n", path)
}

// CreateTable создает талицу с полями и индексирует поле дата
func CreateTable(db *sql.DB) {
	scheduler_table := `CREATE TABLE scheduler (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        "date" CHAR(8) NOT NULL DEFAULT "",
        "title" TEXT NOT NULL DEFAULT "",
        "comment" TEXT NOT NULL DEFAULT "",
        "repeat" VARCHAR(128) NOT NULL DEFAULT "");
		CREATE INDEX scheduler_date ON scheduler (date);`
	query, err := db.Prepare(scheduler_table)
	if err != nil {
		log.Fatal(err)
	}
	query.Exec()
	fmt.Println("Table created successfully!")
}
