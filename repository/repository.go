package repository

import (
	"os"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

const (
	SQLCreateScheduler = `
	CREATE TABLE scheduler (
	    id      INTEGER PRIMARY KEY, 
	    date    CHAR(8) NOT NULL DEFAULT "", 
	    title   TEXT NOT NULL DEFAULT "" CHECK (length(title) < 128),
		comment TEXT NOT NULL DEFAULT "",
		repeat  VARCHAR(128) NOT NULL DEFAULT "" 
	);
	`

	SQLCreateSchedulerIndex = `
	CREATE INDEX scheduler_date_index ON scheduler (date)
	`
)

func NewDB(dbFile string) (*sqlx.DB, error) {
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
