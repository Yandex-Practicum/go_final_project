package db

import (
	"database/sql"
	"fmt"
	"github.com/kmakasheva/todo-list-project/logger"
	"log"
	"os"
	"path"

	_ "modernc.org/sqlite"
)

func PathToDB() (string, error) {
	appPath := os.Getenv("TODO_DBFILE")

	if appPath == "" {
		nowPath, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("error getting current working directory: %w", err)
		}
		appPath = nowPath
	}

	dbFile := path.Join(appPath, "scheduler.db")
	return dbFile, nil
}

func OpenDB(dbFile string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, fmt.Errorf("error while opening db %w", err)
	}
	return db, nil
}

func CreateDB() *sql.DB {
	dbFile, err := PathToDB()
	if err != nil {
		log.Fatalf("Error while initializing db")
	}
	_, err = os.Stat(dbFile)

	var install bool
	if err != nil {
		if os.IsNotExist(err) {
			install = true
		} else {
			logger.Log.Error("error:", logger.Err(err))
			os.Exit(1)
		}
	}

	_, err = os.Create(dbFile)
	if err != nil {
		logger.Log.Error("Error creating database file", logger.Err(err))
		os.Exit(1)
	}

	db, err := OpenDB(dbFile)
	if err != nil {
		logger.Log.Error("error while connection to db", logger.Err(err))
		os.Exit(1)
	}

	_, err = db.Exec(CreateTableQuery)
	if err != nil {
		logger.Log.Error("Error creating table:", logger.Err(err))
		os.Exit(1)
	}
	if install {
		_, err = db.Exec(CreateIndex)
		if err != nil {
			logger.Log.Error("Error creating index by date:", logger.Err(err))
			os.Exit(1)
		}
	}
	return db
}
