package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/PhilippElizarov/go_final_project/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

const dbName string = "scheduler.db"

var dbFile string

const timeTemplate = "20060102"

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal(err.Error())
	}

	var dir string

	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	dir, exists := os.LookupEnv("TODO_DBFILE")
	if exists {
		appPath = dir
	}

	dbFile = filepath.Join(filepath.Dir(appPath), dbName)

	var install bool
	_, err = os.Stat(dbFile)
	if err != nil {
		install = true
	}

	if install {
		file, err := os.Create(dbFile)
		if err != nil {
			log.Fatal(err.Error())
		}
		file.Close()

		sqliteDatabase, _ := sql.Open("sqlite3", dbFile)
		defer sqliteDatabase.Close()
		database.CreateTable(sqliteDatabase)
	}

	router := NewRouter()

	port, exists := os.LookupEnv("TODO_PORT")
	if !exists {
		port = "7540"
	}

	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err.Error())
	}
}
