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
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	install := os.IsNotExist(err)

	conn, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	if install {
		if _, err := conn.Exec(schema); err != nil {
			conn.Close()
			return err
		}
	}

	if err := conn.Ping(); err != nil {
		conn.Close()
		return err
	}

	database = conn
	return nil
}

func Close() error {
	if database == nil {
		return nil
	}

	err := database.Close()
	database = nil
	return err
}
