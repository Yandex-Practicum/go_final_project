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

	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		log.Println("Database file not found, creating new one:", dbFile)
	}

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

	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name='scheduler';")
	if err != nil {
		closeDB()
		return nil, nil, fmt.Errorf("error checking table existence: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		err = createTableAndIndex(db)
		if err != nil {
			closeDB()
			return nil, nil, fmt.Errorf("error creating table and index: %w", err)
		}
	}
	return db, closeDB, nil
}

// The path to the database file.
func pathFileDB() string {
	// Checking the environment variable.
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile != "" {
		return dbFile
	}

	// Get the path to the executable file and create the path to the database.
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("The file path is missing: %v", err)
	}
	dbFile = filepath.Join(filepath.Join(currentDir), "scheduler.db")
	return dbFile
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

func GetTasks(db *sqlx.DB, search string, dateSearch string) ([]task.Task, error) {
	var tasks []task.Task
	var query string
	var args []interface{}

	if dateSearch != "" {
		query = "SELECT * FROM scheduler WHERE date = ? ORDER BY date LIMIT 50"
		args = append(args, dateSearch)
	} else if search != "" {
		searchPattern := "%" + search + "%"
		query = `SELECT * FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT 50`
		args = append(args, searchPattern, searchPattern)
	} else {
		query = `SELECT * FROM scheduler ORDER BY date LIMIT 50`
	}

	err := db.Select(&tasks, query, args...)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func AddTask(db *sqlx.DB, task *task.Task) (int64, error) {
	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)"

	res, err := db.NamedExec(query, task)
	if err != nil {
		return 0, fmt.Errorf("error inserting task: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error getting last insert ID: %w", err)
	}
	return id, nil
}
