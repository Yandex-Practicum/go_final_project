package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

// Глобальная переменная для базы данных
var db *sqlx.DB
var once sync.Once
var mu sync.Mutex

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(100) NOT NULL,
    comment TEXT NOT NULL DEFAULT '',
    repeat VARCHAR(128) NOT NULL DEFAULT ''
);
CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);
`

// Init инициализирует подключение к базе данных и создаёт таблицу
func Init(dbPath string) error {
	mu.Lock()
	defer mu.Unlock()

	if db != nil {
		return nil
	}

	var err error
	db, err = sqlx.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Устанавливаем режим журнала WAL
	_, err = db.Exec("PRAGMA journal_mode=WAL;")
	if err != nil {
		return fmt.Errorf("failed to set journal mode: %w", err)
	}

	// Устанавливаем кодировку UTF-8
	_, err = db.Exec("PRAGMA encoding='UTF-8';")
	if err != nil {
		return fmt.Errorf("failed to set encoding: %w", err)
	}

	// Создаём таблицу, если она не существует
	query := `
        CREATE TABLE IF NOT EXISTS scheduler (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            date TEXT NOT NULL,
            title TEXT NOT NULL,
            comment TEXT,
            repeat TEXT
        );
    `
	_, err = db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	return nil
}

// AddTask добавляет задачу в базу данных и возвращает её ID
func AddTask(task Task) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mu.Lock()
	defer mu.Unlock()

	if db == nil {
		err := Init("../scheduler.db")
		if err != nil {
			return 0, err
		}
	}

	if db == nil {
		return 0, errors.New("database not initialized")
	}

	result, err := db.ExecContext(ctx,
		"INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)",
		task.Date,
		task.Title,
		task.Comment,
		task.Repeat,
	)
	if err != nil {
		log.Printf("Failed to insert task: %v", err)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Failed to get last insert ID: %v", err)
		return 0, err
	}

	return id, nil
}
