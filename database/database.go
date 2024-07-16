package database

import (
	"database/sql"
	_ "embed"
	"os"
	"path/filepath"
)

type DbHelper struct {
	Db *sql.DB
}

var createTableSQL string

var createIndexSQL string

func InitDb() (*DbHelper, error) {
	appPath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")
	_, err = os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		return nil, err
	}

	dbHelper := &DbHelper{Db: db}

	if !install {
		if err := dbHelper.createTables(); err != nil {
			return nil, err
		}
	}

	return dbHelper, nil
}

func (d *DbHelper) createTables() error {
	_, err := d.Db.Exec(createTableSQL)
	if err != nil {
		return err
	}
	_, err = d.Db.Exec(createIndexSQL)
	return err
}
