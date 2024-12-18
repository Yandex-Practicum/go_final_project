package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"final_project/models"

	_ "github.com/mattn/go-sqlite3"
)

func InitDatabase() (*sql.DB, error) {
	appPath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")

	_, err = os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	if install {
		if err := createTable(db); err != nil {
			return nil, err
		}
	}

	return db, nil
}

func GetTasks() ([]models.Task, error) {
	db, err := sql.Open("sqlite3", "scheduler.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT 50")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func AddTask(date, title, comment, repeat string) (string, error) {
	db, err := sql.Open("sqlite3", "scheduler.db")
	if err != nil {
		return "", err
	}
	defer db.Close()

	result, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)", date, title, comment, repeat)
	if err != nil {
		return "", err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return "", err
	}

	return strconv.FormatInt(id, 10), nil
}

func createTable(db *sql.DB) error {
	createTableSQL := `CREATE TABLE scheduler (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        date TEXT NOT NULL,
        title TEXT NOT NULL,
        comment TEXT,
        repeat TEXT
    );`

	statement, err := db.Prepare(createTableSQL)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec()
	if err != nil {
		return err
	}

	createIndexSQL := `CREATE INDEX idx_date ON scheduler (date);`
	indexStatement, err := db.Prepare(createIndexSQL)
	if err != nil {
		return err
	}
	defer indexStatement.Close()

	_, err = indexStatement.Exec()
	if err != nil {
		return err
	}

	return nil
}

func GetTaskByID(id string) (models.Task, error) {
	var task models.Task
	db, err := sql.Open("sqlite3", "scheduler.db")
	if err != nil {
		return task, err
	}
	defer db.Close()

	err = db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return task, fmt.Errorf("задача не найдена")
		}
		return task, err
	}

	return task, nil
}

func UpdateTask(task models.Task) error {
	db, err := sql.Open("sqlite3", "scheduler.db")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?",
		task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("задача не найдена")
		}
		return err
	}

	return nil
}

func DeleteTask(id string) error {
	db, err := sql.Open("sqlite3", "scheduler.db")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM scheduler WHERE id = ?", id)
	return err
}
