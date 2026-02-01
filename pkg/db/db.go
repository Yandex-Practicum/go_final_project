package db

import (
	"database/sql"
	"errors"
	"os"

	"github.com/JKasus/go_final_project/pkg/constants"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

func Init(dbFile string) error {
	_, err := os.Stat(dbFile)

	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return errors.New("error opening database: " + err.Error())
	}

	if _, err = db.Exec(constants.Schema); err != nil {
		return errors.New("schema creating error: " + err.Error())
	}

	return nil
}
