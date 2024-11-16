package main

import (
	"database/sql"

	"log"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var timeFormat = "20060102"

func main() {
	dbCheck()
	var webDir = "./web"
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	http.HandleFunc("/api/nextdate", nextDateHandler)
	http.HandleFunc("/api/task", auth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			taskPost(w, r)
		case http.MethodGet:
			getTask(w, r)
		case http.MethodPut:
			taskPut(w, r)
		case http.MethodDelete:
			taskDelete(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))
	http.HandleFunc("/api/signin", sign)
	http.HandleFunc("/api/tasks", auth(getTasks))
	http.HandleFunc("/api/task/done", auth(doTask))

	http.ListenAndServe(":"+os.Getenv("TODO_PORT"), nil)

	log.Println("Starting server on :" + os.Getenv("TODO_PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("TODO_PORT"), nil))
}

func dbCheck() {
	var install bool
	var dbFile string
	if os.Getenv("TODO_DBFILE") == "" {
		appPath, _ := os.Executable()
		_, err := os.Stat(filepath.Join(filepath.Dir(appPath), "scheduler.db"))
		if err != nil {
			install = true
			dbFile = filepath.Join(filepath.Dir(appPath), "scheduler.db")
		}
	} else {
		_, err := os.Stat(os.Getenv("TODO_DBFILE"))
		if err != nil {
			install = true
			appPath, _ := os.Executable()
			dbFile = filepath.Join(filepath.Dir(appPath), "scheduler.db")
		} else {
			dbFile = os.Getenv("TODO_DBFILE")
		}
	}
	if install {
		db, err := sql.Open("sqlite3", dbFile)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		_, err = db.Exec(`CREATE TABLE IF NOT EXISTS scheduler(
		id INTEGER PRIMARY KEY,
		date TEXT,
		title TEXT,
		comment TEXT,
		repeat TEXT)`)
		if err != nil {
			log.Fatal(err)
		}
		_, err = db.Exec("CREATE INDEX date_index ON scheduler(date)")
		if err != nil {
			log.Println(err)
		}
	}
}
