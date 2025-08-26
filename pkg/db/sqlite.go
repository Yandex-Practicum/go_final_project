package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func ConnectDB(dbFile string) (*sql.DB, error) {
	return sql.Open("sqlite3", dbFile)
}
