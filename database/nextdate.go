package storage

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"

	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type DbConnection struct {
	Db *sqlx.DB
}

type scheduler struct {
	ID      int
	Date    time.Duration
	Title   string
	Comment string
	Repeat  string
}

func NewDb() (*DbConnection, error) {
	file, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Println(file)

	dbFile := filepath.Join(file, os.Getenv("DB_NAME"))

	_, err = os.Stat(dbFile)
	var install bool
	if err != nil {
		install = true
	}
	db, err := sqlx.Open(os.Getenv("DB_DRIVER"), dbFile)
	if err != nil {
		log.Printf("DbConnection: can't open db: %s\n", err)
	}

	if install {
		_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date СHAR(8) NOT NULL,
			title VARCHAR NOT NULL,
			comment TEXT,
			repeat VARCHAR(128)
		);
		CREATE INDEX idx_scheduler_date ON scheduler (date);
		`)
		if err != nil {
			log.Printf("DbConnection: can't create db: %s\n", err)
		}
	}

	log.Printf("Database connected. Database file: %s", dbFile)
	return &DbConnection{Db: db}, nil
}

// TODO: close db
func (db *DbConnection) Close() {
	db.Db.Close()

}
