package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"go_final_project/config"

	_ "modernc.org/sqlite"
)

const dbDriver = "sqlite"

var dbPath = config.DbPath()

type Storage struct {
	db *sql.DB
}

func New() (*Storage, error) {
	s := &Storage{}

	err := s.initDB()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Storage) initDB() error {
	var err error

	err = createDB(dbPath)
	if err != nil {
		return err
	}

	db, err := sql.Open(dbDriver, dbPath)
	if err != nil {
		log.Println("Can't open database")
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}

	if err = createNewTable(db); err != nil {
		return err
	}

	s.db = db

	return nil
}

// функция Close закрывает соединение с базой данных
func (s *Storage) Close() error {
	if s.db != nil {
		if err := s.db.Close(); err != nil {
			return fmt.Errorf("can`t close database connection: %w", err)
		}
	}
	return nil
}

func createDB(dbPath string) error {
	_, err := os.Create(dbPath)
	if err != nil {
		return fmt.Errorf("can't create storage by func createDB: %w", err)
	}

	log.Println("The database file has been created")

	return nil
}

func createNewTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS scheduler (
		   id INTEGER PRIMARY KEY AUTOINCREMENT,
		   date VARCHAR(8) NOT NULL,
		   title TEXT NOT NULL,
		   comment TEXT DEFAULT "",
		   repeat VARCHAR(128) NOT NULL
   		);
	
   		CREATE INDEX scheduler_date ON scheduler (date);
   `)

	if err != nil {
		return fmt.Errorf("can't create new table by func createNewTable: %w", err)
	}

	return nil
}
