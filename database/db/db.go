package db

import (
	"database/sql"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

const (
	maxIdleConns    = 2
	maxOpenConns    = 5
	connMaxIdleTime = time.Minute * 5
)

var (
	DateFormat string
)

// DBExists проверяет существует ли файл переданный аргументом
func DbExists(dbFile string) bool {
	_, err := os.Stat(dbFile)
	var exists bool
	if err == nil {
		exists = true
	}
	return exists

}

// StartDB открывает базу данных указанную в .env файле, добавляет её в структуру DBHandler.
func StartDB() (Storage, error) {
	dbFile := os.Getenv("TODO_DBFILE")
	DateFormat = os.Getenv("TODO_DATEFORMAT")

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return Storage{}, err
	}
	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxIdleTime(connMaxIdleTime)
	dbStorage := Storage{}
	dbStorage.db = db

	return dbStorage, nil
}

// CloseDB закрывает подключение к базе данных.
func (dbHandl *Storage) CloseDB() error {
	db := dbHandl.db
	err := db.Close()
	if err != nil {
		return err
	}
	return nil

}

// InstallDB создаёт файл для базы данных с названием, указаным в .env,
// отправляет SQL запрос на создание таблицы из файла schema.sql.
// Возвращает ошибку в случае неудачи.
func InstallDB() error {
	dbFile := os.Getenv("TODO_DBFILE")
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return err
	}
	defer db.Close()
	installQuery, err := os.ReadFile("internal/db/schema.sql")
	if err != nil {
		return err
	}
	_, err = db.Exec(string(installQuery))
	if err != nil {
		return err
	}

	return nil
}
