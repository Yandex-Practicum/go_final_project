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

// Инициализация базы данных.
func InitDB() (*sqlx.DB, func(), error) {
	dbFile := pathFileDB()
	log.Println("dbFile:", dbFile)

	db, err := sqlx.Open("sqlite3", dbFile)
	if err != nil {
		return nil, nil, fmt.Errorf("error to open database: %w", err)
	}

	closeDB := func() {
		err := db.Close()
		if err != nil {
			log.Printf("error to close database: %v", err)
		}
	}

	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		log.Println("Database file doesn't exist, creating...")
		err = createTableAndIndex(db)
		if err != nil {
			closeDB()
			return nil, nil, fmt.Errorf("createTableAndIndex: %w", err)
		}
		log.Println("Database created successfully")
	}
	return db, closeDB, nil
}

// Путь к файлу базы данных.
func pathFileDB() string {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile != "" {
		log.Printf("The database file from the environment variable TODO_DBFILE: %s", dbFile)
		return dbFile
	}

	appPath, err := os.Executable()
	if err != nil {
		log.Fatalf("The file path is missing: %v", err)
	}
	dbFile = filepath.Join(filepath.Dir(appPath), "scheduler.db")
	log.Printf("File database default: %s", dbFile)
	return dbFile
}

// Создание таблицы scheduler и индекса по полю date
func createTableAndIndex(db *sqlx.DB) error {
	log.Println("Creating table scheduler")
	// Создание таблицы по полю date
	createTable := `CREATE TABLE scheduler (id INTEGER PRIMARY KEY AUTOINCREMENT, date TEXT NOT NULL, title TEXT NOT NULL, comment TEXT, repeat TEXT);`
	// Создание индекса по полю ate
	createIndex := `CREATE INDEX idx_date ON scheduler (date);`
	_, err := db.Exec(createTable)
	if err != nil {
		log.Printf("error to create table: %v", err)
		return fmt.Errorf("error to create table: %w", err)
	}
	log.Println("Table scheduler created successfully.")

	log.Println("Creating indexes")
	_, err = db.Exec(createIndex)
	if err != nil {
		log.Printf("error to create indexes: %v", err)
		return fmt.Errorf("error to create index: %w", err)
	}
	log.Println("Indexes created successfully.")
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
