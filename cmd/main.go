package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"pwd/handlers"

	"github.com/go-chi/chi/v5"
	_ "modernc.org/sqlite"
)

func main() {
	webDir := "./web"

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
	// если install равен true, после открытия БД требуется выполнить
	// sql-запрос с CREATE TABLE и CREATE INDEX
	if install {

		db, err := sql.Open("sqlite", dbFile)
		if err != nil {
			fmt.Printf("Ошибка при попытке соединения с базой данный: %s\n", err.Error())
			return
		}
		defer db.Close()

		createDb := `CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(256) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat VARCHAR(128) NOT NULL DEFAULT "");

	CREATE INDEX tasks_rules ON scheduler (repeat);`

		_, err = db.Exec(createDb)
		if err != nil {
			fmt.Printf("Ошибка создания таблицы: %s\n", err.Error())
			return
		}

	}

	r := chi.NewRouter()
	r.Handle("/", http.FileServer(http.Dir(webDir)))
	r.Get("/api/nextdate", handlers.NextDateHandler)
	err = http.ListenAndServe(":7540", r)
	if err != nil {
		log.Fatal(err)
	}
}
