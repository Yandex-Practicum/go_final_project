package database

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// CheckDatabase проверка существования БД
func CheckDatabase() (*sql.DB, error) {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	_, err := os.Stat(dbFile)
	if os.IsNotExist(err) {
		file, err := os.Create(dbFile)
		if err != nil {
			return nil, err
		}
		file.Close()
	}

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	createTable := `CREATE TABLE IF NOT EXISTS scheduler (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
		date CHAR(8) NOT NULL DEFAULT "",
		title VARCHAR(128) NOT NULL DEFAULT "",
		comment TEXT NOT NULL DEFAULT "",
		repeat VARCHAR(128) NOT NULL DEFAULT ""
    );
	CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);`

	_, err = db.Exec(createTable)
	if err != nil {
		return nil, err
	}
	return db, nil
}
