package db

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

const Schema string = `CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(20) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat VARCHAR(128) NOT NULL DEFAULT ""
);

CREATE INDEX IF NOT EXISTS idx_date ON scheduler (date);`

func Init(dbFile string) error {

	dbase, err := sql.Open("sqlite", dbFile)
	if err != nil {
		log.Println(err)
		return err
	}
	defer dbase.Close()

	_, err = dbase.Exec(Schema)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
