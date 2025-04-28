package db

import (
	"database/sql"
	_ "modernc.org/sqlite"
	"os"
)

// Константа для создания БД
const schema = `
CREATE TABLE scheduler (
id      INTEGER PRIMARY KEY AUTOINCREMENT,
date    CHAR(8)     NOT NULL DEFAULT '',
title   VARCHAR(255) NOT NULL DEFAULT '',
comment TEXT,
repeat  VARCHAR(128) NOT NULL DEFAULT ''
);
CREATE INDEX idx_scheduler_date ON scheduler(date);
`

var db *sql.DB

// Init инициализирует БД
func Init(dbFile string) error {
	// Проверка на существование БД
	_, err := os.Stat(dbFile)
	var install bool
	if err != nil {
		install = true
	}

	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}
	// При отсутствии БД, создаем
	if install {
		_, err := db.Exec(schema)
		if err != nil {
			return err
		}
	}
	return nil
}
