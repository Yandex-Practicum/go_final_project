package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	appPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	println(appPath)

	dbFile := filepath.Join(appPath, "scheduler.db")
	_, err = os.Stat(dbFile)
	println(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	if install { //если файл БД не существует, то создаем его

		db, err := sql.Open("sqlite3", dbFile)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		_, err = db.Exec("CREATE TABLE IF NOT EXISTS scheduler (id INTEGER PRIMARY KEY, date TEXT, title TEXT, comment TEXT, repeat TEXT(128))") // создаем таблицу с данными
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec("CREATE INDEX IF NOT EXISTS idx_tasks_date ON scheduler (date)")
		if err != nil {
			log.Fatal(err)
		}

	}

	webDir := "./web"
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	err = http.ListenAndServe(":7541", nil)
	if err != nil {
		panic(err)
	}

}
