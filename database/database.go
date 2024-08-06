// package database provides database setup for working on scheduler tasks
package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// createTableIfNotExists creates table if not exists
func createTableIfNotExists(db *sql.DB) error {
	// table creation query
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS scheduler (
	    id INTEGER PRIMARY KEY AUTOINCREMENT, 
	    date CHAR(8) NOT NULL DEFAULT "", 
	    title TEXT NOT NULL DEFAULT "",
		comment TEXT,
		repeat VARCHAR(128)
	);
	`
	// table creation
	_, err := db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("error creating table: %w", err)
	}
	// index creation query
	createIndexSQL := "CREATE INDEX IF NOT EXISTS index_date ON scheduler (date);"

	//index creation
	_, err = db.Exec(createIndexSQL)
	if err != nil {
		return fmt.Errorf("error creating column index: %w", err)
	}

	return nil
}

// SetupDatabase configure the DB connection
func SetupDatabase() (*sql.DB, error) {
	// get enviroment variable of db
	dbFile := os.Getenv("TODO_DBFILE")

	// get db path if variable is empty
	if dbFile == "" {
		appPath, err := os.Executable()
		if err != nil {
			return nil, fmt.Errorf("error getting executable path: %w", err)
		}
		dbFile = filepath.Join(filepath.Dir(appPath), "scheduler.db")
	}

	// check for file existence
	_, err := os.Stat(dbFile)
	var install bool
	if err != nil {
		if os.IsNotExist(err) {
			install = true
		} else {
			return nil, fmt.Errorf("error checking database file: %w", err)
		}
	}

	// establish connection with db
	dbCon, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	if install {
		if err := createTableIfNotExists(dbCon); err != nil {
			return nil, fmt.Errorf("failed to set up database: %w", err)
		}
	}

	return dbCon, nil
}
