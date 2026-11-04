package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

var database *sql.DB

const schema = `
CREATE TABLE scheduler (
id INTEGER PRIMARY KEY AUTOINCREMENT,
date CHAR(8) NOT NULL DEFAULT "",
title VARCHAR(256) NOT NULL DEFAULT "",
comment TEXT NOT NULL DEFAULT "",
repeat VARCHAR(128) NOT NULL DEFAULT ""
);

CREATE INDEX scheduler_date ON scheduler (date);
`

func Init(dbFile string) error {
	_, err := os.Stat(dbFile)
	install := os.IsNotExist(err)

	var openErr error
	database, openErr = sql.Open("sqlite", dbFile)
	if openErr != nil {
		return openErr
	}

	if install {
		if _, err := database.Exec(schema); err != nil {
			return err
		}
	}

	return database.Ping()
}
