package base

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func CreateDB(envDBFILE string) {
	var appPath string
	var err error
	if envDBFILE != "" {
		appPath = envDBFILE
	} else {
		appPath, err = os.Getwd() //не смогла реализовать через os.Executable()
		if err != nil {
			log.Fatal(err)
		}
	}

	dbFile := filepath.Join(appPath, "scheduler.db")
	_, err = os.Stat(dbFile)

	fmt.Println(err)

	var install bool
	if os.IsNotExist(err) {
		install = true
		fmt.Println("db не найдена")
	} else if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if install {

		// Создание базы данных
		_, err = db.Exec(`CREATE TABLE IF NOT EXISTS scheduler(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date INTEGER,
			title TEXT NOT NULL DEFAULT "",
			comment TEXT,
			repeat VARCHAR(128) NOT NULL DEFAULT "");`)
		if err != nil {
			log.Fatal(err)
		}
		_, err = db.Exec(`CREATE INDEX IF NOT EXISTS scheduler_date ON scheduler (date);`)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("База данных успешно создана!")
	}
}
