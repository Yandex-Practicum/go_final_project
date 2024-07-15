package database

import (
	"os"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

type Database struct {
	Db *sqlx.DB
}

var DB Database

func ConnectDB() error {
	dbFile := os.Getenv("TODO_DBFILE")
	var install bool
	_, err := os.Stat(dbFile)
	if err != nil {
		install = true
	}
	// если install равен true, после открытия БД требуется выполнить
	// sql-запрос с CREATE TABLE и CREATE INDEX
	if install {
		_, err = os.Create(dbFile)
		if err != nil {
			return err
		}
	}

	db, err := sqlx.Connect("sqlite", dbFile)
	if err != nil {
		return err
	}

	if install {
		err = createTable(db)
		if err != nil {
			return err
		}
	}

	DB = Database{
		Db: db,
	}

	return nil
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
