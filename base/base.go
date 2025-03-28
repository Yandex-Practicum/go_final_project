package base

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// Создание базы данных 2ой шаг
func CreateDB(envDBFILE string) (db *sql.DB, err error) {
	var appPath string
	if envDBFILE != "" {
		appPath = envDBFILE
	} else {
		appPath, err = os.Getwd() //не смогла реализовать через os.Executable()
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
	}

	var install bool
	dbFile := filepath.Join(appPath, "scheduler.db")
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		install = true
		fmt.Println("db не найдена, создаём новую")
	}

	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	if install {

		// Создание базы данных
		_, err = db.Exec(`CREATE TABLE IF NOT EXISTS scheduler(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date TEXT NOT NULL,
			title TEXT NOT NULL DEFAULT "",
			comment TEXT,
			repeat VARCHAR(128) NOT NULL DEFAULT "");`)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		_, err = db.Exec(`CREATE INDEX IF NOT EXISTS scheduler_date ON scheduler (date);`)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		fmt.Println("База данных успешно создана!")
	}
	return db, nil
}
