package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"go_final-project/internal/task"
	"log"
	"os"
	"path/filepath"
)

// Initializing the database.
func InitDB() (*sqlx.DB, func(), error) {
	dbFile := pathFileDB()

	// Opening a database connection.
	db, err := sqlx.Open("sqlite3", dbFile)
	if err != nil {
		return nil, nil, fmt.Errorf("error to open database: %w", err)
	}

	// A function for closing the database.
	closeDB := func() {
		err := db.Close()
		if err != nil {
			log.Printf("error to close database: %v", err)
		}
	}

	// Checking if the scheduler table exists.
	var exists int
	err = db.Get(&exists, "SELECT count(*) FROM sqlite_master WHERE type='table' AND name='scheduler'")
	if err != nil || exists == 0 {
		log.Println("Table 'scheduler' not found. Creating new database structure.")
		err = createTableAndIndex(db)
		if err != nil {
			closeDB()
			return nil, nil, fmt.Errorf("error to create table: %w", err)
		}
	}
	log.Println("Database initialized successfully")
	return db, closeDB, nil
}

// The path to the database file.
func pathFileDB() string {
	// Checking the environment variable.
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile != "" {
		return dbFile
	}

	// We get the path to the executable file and create the path to the database.
	appPath, err := os.Executable()
	if err != nil {
		log.Fatalf("The file path is missing: %v", err)
	}
	return filepath.Join(filepath.Dir(appPath), "scheduler.db")
}

// Creating a scheduler table and an index for the date field
func createTableAndIndex(db *sqlx.DB) error {
	// Создание таблицы по полю date
	createTable := `CREATE TABLE scheduler (id INTEGER PRIMARY KEY AUTOINCREMENT, date TEXT NOT NULL, title TEXT NOT NULL, comment TEXT, repeat TEXT);`
	// Создание индекса по полю ate
	createIndex := `CREATE INDEX idx_date ON scheduler (date);`
	_, err := db.Exec(createTable)
	if err != nil {
		return fmt.Errorf("error to create table: %w", err)
	}
	_, err = db.Exec(createIndex)
	if err != nil {
		return fmt.Errorf("error to create index: %w", err)
	}

	log.Println("Table 'scheduler' and index 'idx_date' created successfully")
	return nil
}

func GetTasks(db *sqlx.DB) ([]task.Task, error) {
	var tasks []task.Task
	err := db.Select(&tasks, "SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date")
	if err != nil {
		return nil, fmt.Errorf("error to get tasks: %w", err)
	}
	return tasks, nil
}

func AddTask(db *sqlx.DB, task *task.Task) (int64, error) {
	res, err := db.NamedExec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)", task)
	if err != nil {
		return 0, fmt.Errorf("error to add task: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error to get ID: %w", err)
	}
	return id, nil
}
