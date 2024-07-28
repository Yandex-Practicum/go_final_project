package storage

import (
	"database/sql"
	"fmt"
	"go_final_project/internal/config"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

type Storage struct {
	db *sql.DB
}

var dbPath = config.StoragePath
var dbDriver = config.Driverdb

func New() (*Storage, error) {
	storage := &Storage{}

	err := storage.initDB()

	if err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *Storage) initDB() error {
	var err error

	err = createDB(dbPath)
	if err != nil {
		return err
	}

	db, err := sql.Open(dbDriver, dbPath)

	if err != nil {
		log.Println("Не удалось открыть базу данных")
		return err
	}

	err = createTableAndIdx(db)

	if err != nil {
		return err
	}

	s.db = db

	return nil
}

func createDB(path string) error {
	_, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("не удалось создать хранилище %w", err)
	}
	return nil
}

func createTableAndIdx(db *sql.DB) error {
	_, err := db.Exec(`
    CREATE TABLE IF NOT EXISTS scheduler (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      date CHAR(8) NOT NULL,
      title VARCHAR(128) NOT NULL DEFAULT '',
      comment TEXT DEFAULT '',
      repeat VARCHAR(128) NOT NULL
      );
    CREATE INDEX idx_date ON scheduler (date);
  `)
	if err != nil {
		return fmt.Errorf("не удалось создать новую таблицу %w", err)
	}

	return nil
}
