package db

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"go_final-project/internal/task"
	"log"
	"os"
	"path/filepath"
)

// InitDB initializing the database.
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

// pathFileDB returns the path to the SQLite database file, either from an environment variable or default.
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

// createTableAndIndex initializes the "scheduler" table and creates an index on the "date" field.
func createTableAndIndex(db *sqlx.DB) error {
	// createTable defines the schema for the "scheduler" table,
	// where tasks are stored with an ID, date, title, optional comment, and repeat rule.
	createTable := `CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date TEXT NOT NULL,
    title TEXT NOT NULL,
    comment TEXT,
    repeat TEXT
                       );`
	// createIndex improves query performance by creating an index on the "date" field.
	// This allows for faster searches and sorting of tasks by their due date.
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

// GetTasks retrieves a list of tasks from the database, optionally filtered by date or search query.
func GetTasks(db *sqlx.DB, search string, dateSearch string) ([]task.Task, error) {
	var tasks []task.Task
	var query string
	var args []interface{}

	if dateSearch != "" {
		query = `SELECT * FROM scheduler WHERE date = ? ORDER BY date LIMIT 50`
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

// GetTaskByID retrieves a task from the "scheduler" table by its unique ID.
func GetTaskByID(db *sqlx.DB, id int64) (*task.Task, error) {
	var t task.Task
	query := "SELECT * FROM scheduler WHERE id = ?"
	err := db.Get(&t, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}
	return &t, nil
}

// AddTask inserts a new task into the "scheduler" table.
// It returns the unique ID of the newly created task or an error if the insertion fails.
func AddTask(db *sqlx.DB, task *task.Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)`

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

// UpdateTask updates an existing task in the "scheduler" table by its ID.
// It modifies the task's date, title, comment, and repeat rule.
// Returns an error if the update fails or if the task does not exist.
func UpdateTask(db *sqlx.DB, t *task.Task) error {
	query := `UPDATE scheduler
	SET date=:date, title=:title, comment=:comment, repeat=:repeat
	WHERE id = :id`

	res, err := db.NamedExec(query, t)
	if err != nil {
		return fmt.Errorf("error updating task: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking changes: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("issue not found")
	}
	return nil
}

// DeleteTask removes a task from the "scheduler" table by its ID.
// If the task does not exist, it returns an error.
func DeleteTask(db *sqlx.DB, id int64) error {
	query := "DELETE FROM scheduler WHERE id = ?"
	res, err := db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("issue not found")
	}
	return nil
}
