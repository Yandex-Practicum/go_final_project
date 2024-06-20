package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

const defaultPort = "7540"

func main() {
	port := getPort()
	webDir := "web"

	checkDb()

	fs := http.FileServer(http.Dir(webDir))
	http.Handle("/", fs)

	fmt.Println("Запуск сервера на 7540....")
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}
func getPort() string {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = defaultPort
	}
	return port
}

func checkDb() {
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

	if !install {
		db, err := sql.Open("sqlite", "scheduler.db")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer db.Close()

		queryCreate, err := db.Query(`
		CREATE TABLE scheduler 
		(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date CHAR(8) NOT NULL DEFAULT "",
		title TEXT NOT NULL DEFAULT "",
		comment TEXT NOT NULL DEFAULT "",
		repeat VARCHAR(128) NOT NULL DEFAULT ""
		);
		CREATE INDEX date_index ON scheduler (date);`)

		if err != nil {
			fmt.Println(err)
			return
		}

		defer queryCreate.Close()
		/*db.SetMaxIdleConns(2)
		db.SetMaxOpenConns(5)
		db.SetConnMaxIdleTime(time.Minute * 5)
		db.SetConnMaxLifetime(time.Hour)*/
	}
}
