package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

var conn *sql.DB

const schema = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT '',
    title VARCHAR(255) NOT NULL DEFAULT '',
    comment TEXT NOT NULL DEFAULT '',
    repeat VARCHAR(128) NOT NULL DEFAULT ''
    );
CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);
`

func Init(dbFile string) error {
	install := false
	if _, err := os.Stat(dbFile); err != nil {
		install = true
	}

	c, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	err = c.Ping()
	if err != nil {
		_ = c.Close()
		return err
	}

	if install {
		_, err = c.Exec(schema)
		if err != nil {
			_ = c.Close()
			return err
		}
	}

	conn = c
	return nil
}

func DB() *sql.DB { return conn }

func Close() error {
	if conn != nil {
		err := conn.Close()
		conn = nil
		return err
	}
	return nil
}
