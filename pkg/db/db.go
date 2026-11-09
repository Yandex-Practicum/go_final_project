package db

import (
	"database/sql"
	_ "modernc.org/sqlite"
	"os"
)

const schema = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT '',
    title VARCHAR(256) NOT NULL,
    comment TEXT,
    repeat VARCHAR(128)
);

CREATE INDEX idx_scheduler_date ON scheduler(date);
`

func Init(dbFile string) (*sql.DB, error) {
	_, err := os.Stat(dbFile)

	install := false
	if os.IsNotExist(err) {
		install = true
	} else if err != nil {
		return nil, err
	}

	database, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, err
	}

	if err = database.Ping(); err != nil {
		database.Close()
		return nil, err
	}

	if install {
		_, err = database.Exec(schema)
		if err != nil {
			database.Close()
			return nil, err
		}
	}

	return database, nil
}
