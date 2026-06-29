package db

import (
	"database/sql"
	"errors"
	"os"

	_ "modernc.org/sqlite"
)

const schema = `CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT '',
    title VARCHAR(128) NOT NULL DEFAULT '',
    comment TEXT NOT NULL DEFAULT '',
    repeat VARCHAR(128) NOT NULL DEFAULT ''
);`

const indexSchema = `CREATE INDEX idxDate ON scheduler(date);`

var Db *sql.DB

func Init(dbFile string) error {
	var install bool
	_, err := os.Stat(dbFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			install = true
		} else {
			return err
		}
	}
	Db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}
	if install {
		_, err = Db.Exec(schema)
		if err != nil {
			Db.Close()
			return err
		}
		_, err = Db.Exec(indexSchema)
		if err != nil {
			Db.Close()
			return err
		}
	}
	return nil
}
