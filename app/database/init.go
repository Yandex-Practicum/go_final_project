package database

import (
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() {

	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "./scheduler.db"
	}

	// check if file exists
	_, err := os.Stat(dbFile)
	install := os.IsNotExist(err)

	// open database connection
	DB, err = sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatalf("Error creating database: %v", err)
	}

	// create a table if sqlite file is absent
	if install {
		log.Println("Creating a new database")
		if err := createDB(DB); err != nil {
			log.Fatalf("Error creating database: %v", err)
		}
	}
}

func CloseDB() {
	if DB != nil {
		err := DB.Close()
		if err != nil {
			log.Printf("Error closing database: %v", err)
		} else {
			log.Println("Database connection closed")
		}
	}
}

func createDB(db *sql.DB) error {
	query := ` CREATE TABLE IF NOT EXISTS scheduler (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				date VARCHAR(255) NOT NULL,
				title TEXT NOT NULL,
				comment TEXT,
				repeat VARCHAR(128) NOT NULL);
				CREATE INDEX IF NOT EXISTS idx_date ON scheduler (date);`
	_, err := db.Exec(query)
	return err
}
