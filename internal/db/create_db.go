package db

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var Install bool

func CreateDB() *sql.DB {

	var DBFile string
	if os.Getenv("TODO_DBFILE") == "" {
		appPath, err := os.Executable()
		if err != nil {
			log.Fatal(err)
		}

		if _, err := os.Stat(filepath.Join(filepath.Dir(appPath), "scheduler.db")); err != nil {
			Install = true
		}
	} else {
		path := os.Getenv("TODO_DBFILE")
		_, err := os.Stat(path)
		DBFile = path
		if err != nil {
			Install = true
		}
	}

	db, err := sql.Open("sqlite3", DBFile)
	if err != nil {
		log.Fatal("ошибка создания БД ", err)
	}

	// err = db.Ping()
	// if err != nil {
	// 	log.Fatal("ошибка подключения к базе данных ", err)
	// }
	return db

}
