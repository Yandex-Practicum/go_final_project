package sqliterepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"

	"todo/internal/entity"

	"github.com/jmoiron/sqlx"
	// Importing the sqlite3 driver package.
	// This blank import is necessary to register the driver with the database/sql package.
	_ "github.com/mattn/go-sqlite3"
)

// SQLiteRepository is a repository implementation that uses SQLite as the underlying database.
type SQLiteRepository struct {
	db *sqlx.DB
}

// New creates a new instance of SQLiteRepository.
// It initializes the SQLite database file,
// connects to the database, and creates the necessary schema if the file is newly created.
// It returns the initialized SQLiteRepository instance or an error if any.
func New(dbPath string) (*SQLiteRepository, error) {
	fileExists, err := initDBFile(dbPath)
	if err != nil {
		return nil, fmt.Errorf("initDBFile: %w", err)
	}

	db, err := sqlx.Connect("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("sqlx.Connect: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("db.Ping: %w", err)
	}

	if !fileExists {
		if err = createDBSchema(db); err != nil {
			return nil, fmt.Errorf("createDBSchema: %w", err)
		}
	}

	return &SQLiteRepository{
		db: db,
	}, nil
}

// Close closes the connection to the SQLite database.
// It returns an error if any.
func (sr *SQLiteRepository) Close() error {
	if err := sr.db.Close(); err != nil {
		return fmt.Errorf("db.Close: %w", err)
	}

	return nil
}

// AddTask adds a new task to the SQLite database.
// It takes a context, representing the execution context, and a task, representing the task to be added.
// It returns the ID of the newly added task or an error if any.
func (sr *SQLiteRepository) AddTask(ctx context.Context, task *entity.Task) (int64, error) {
	result, err := sr.db.NamedExecContext(ctx, `
		INSERT INTO scheduler (title, comment, date, repeat)
		VALUES (:title, :comment, :date, :repeat)
	`, task)
	if err != nil {
		return 0, fmt.Errorf("db.Exec: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("result.LastInsertId: %w", err)
	}

	return id, nil
}

// GetTask retrieves a task from the SQLite database by its ID.
// It takes a context, representing the execution context, and an ID, representing the ID of the task to retrieve.
// It returns the retrieved task or an error if any.
func (sr *SQLiteRepository) GetTask(ctx context.Context, id int64) (*entity.Task, error) {
	task := &entity.Task{}
	err := sr.db.GetContext(ctx, task, `
		SELECT id, title, comment, date, repeat
		FROM scheduler
		WHERE id = ?
	`, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.ErrTaskNotFound
		}

		return nil, fmt.Errorf("db.GetContext: %w", err)
	}

	return task, nil
}

// GetTasksByDate retrieves tasks from the SQLite database by their date.
// It takes a context, representing the execution context, and a date, representing the date of the tasks to retrieve.
// It returns the retrieved tasks or an error if any.
func (sr *SQLiteRepository) GetTasksByDate(ctx context.Context, date string) ([]*entity.Task, error) {
	tasks := []*entity.Task{}
	err := sr.db.SelectContext(ctx, &tasks, `
		SELECT id, title, comment, date, repeat
		FROM scheduler
		WHERE date = ?
	`, date)
	if err != nil {
		return nil, fmt.Errorf("db.SelectContext: %w", err)
	}

	return tasks, nil
}

// GetTasksByQuery retrieves tasks from the SQLite database by a search query.
// It takes a context, representing the execution context, and a query, representing the search query.
// It returns the retrieved tasks or an error if any.
func (sr *SQLiteRepository) GetTasksByQuery(ctx context.Context, query string) ([]*entity.Task, error) {
	tasks := []*entity.Task{}
	likeExp := fmt.Sprintf("%%%s%%", strings.ToUpper(query))
	err := sr.db.SelectContext(ctx, &tasks, `
	SELECT id, title, comment, date, repeat 
	FROM scheduler 
	WHERE UPPER(title) LIKE ? OR UPPER(comment) LIKE ? Order by date ASC
	`,
		likeExp, likeExp)
	if err != nil {
		return nil, fmt.Errorf("db.SelectContext: %w", err)
	}

	return tasks, nil
}

// GetTasks retrieves tasks from the SQLite database.
// It takes a context, representing the execution context,
// and a limit, representing the maximum number of tasks to retrieve.
// It returns the retrieved tasks or an error if any.
func (sr *SQLiteRepository) GetTasks(ctx context.Context, limit int) ([]*entity.Task, error) {
	tasks := []*entity.Task{}
	err := sr.db.SelectContext(ctx, &tasks, `
		SELECT id, title, comment, date, repeat
		FROM scheduler
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, fmt.Errorf("db.SelectContext: %w", err)
	}

	return tasks, nil
}

// UpdateTask updates a task in the SQLite database.
// It takes a context, representing the execution context, and a task, representing the updated task.
// It returns an error if any.
func (sr *SQLiteRepository) UpdateTask(ctx context.Context, task *entity.Task) error {
	_, err := sr.db.NamedExecContext(ctx, `
		UPDATE scheduler
		SET title = :title, comment = :comment, date = :date, repeat = :repeat
		WHERE id = :id
	`, task)
	if err != nil {
		return fmt.Errorf("db.ExecContext: %w", err)
	}

	return nil
}

// DeleteTask deletes a task from the SQLite database by its ID.
// It takes a context, representing the execution context, and an ID, representing the ID of the task to delete.
// It returns an error if any.
func (sr *SQLiteRepository) DeleteTask(ctx context.Context, id int64) error {
	_, err := sr.db.ExecContext(ctx, `
		DELETE FROM scheduler
		WHERE id = ?
	`, id)
	if err != nil {
		return fmt.Errorf("db.ExecContext: %w", err)
	}

	return nil
}

// initDBFile initializes the SQLite database file.
// It takes a path, representing the path to the database file.
// It returns a boolean indicating whether the file already exists and an error if any.
func initDBFile(path string) (bool, error) {
	var err error

	_, err = os.Stat(path)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		if _, err = os.Create(path); err != nil {
			return false, fmt.Errorf("os.Create: %w", err)
		}
		return false, nil
	}

	return true, nil
}

// createDBSchema creates the necessary schema in the SQLite database.
// It takes a db, representing the database connection.
// It returns an error if any.
func createDBSchema(db *sqlx.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			comment TEXT,
			date VARCHAR(8) NOT NULL,
			repeat VARCHAR(128)
		);
	`)
	if err != nil {
		return fmt.Errorf("db.Exec: %w", err)
	}

	_, err = db.Exec(`CREATE INDEX idx_scheduler_date ON scheduler (date);`)
	if err != nil {
		return fmt.Errorf("db.Exec: %w", err)
	}

	return nil
}
