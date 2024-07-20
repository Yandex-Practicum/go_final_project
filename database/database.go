package database

import (
	"os"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

func NewDB() (*sqlx.DB, error) {
	dbFile := os.Getenv("TODO_DBFILE")
	var install bool
	_, err := os.Stat(dbFile)
	if err != nil {
		install = true
	}

	if install {
		_, err = os.Create(dbFile)
		if err != nil {
			return nil, err
		}
	}

	db, err := sqlx.Connect("sqlite", dbFile)
	if err != nil {
		return nil, err
	}

	if install {
		err = createTable(db)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}

func createTable(db *sqlx.DB) error {
	_, err := db.Exec(SQLCreateScheduler)
	if err != nil {
		return err
	}

	_, err = db.Exec(SQLCreateSchedulerIndex)
	if err != nil {
		return err
	}

	return nil
}
