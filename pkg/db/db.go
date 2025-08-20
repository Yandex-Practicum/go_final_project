package db

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "modernc.org/sqlite" // Импортируем SQLite-драйвер
)

var DB *sql.DB

const schema = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT '',
    title VARCHAR(256) NOT NULL DEFAULT '',
    comment TEXT NOT NULL DEFAULT '',
    repeat VARCHAR(128) NOT NULL DEFAULT ''
);

CREATE INDEX idx_scheduler_date ON scheduler(date);
`

func Init(dbFile string) error {
	_, err := os.Stat(dbFile)
	install := os.IsNotExist(err)
	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("ошибка при открытии базы данных: %w", err)
	}

	if install {
		_, err := DB.Exec(schema)
		if err != nil {
			return fmt.Errorf("ошибка при открытии базы данных: %w", err)
		}
	}
	return nil
}

func SearchTasks(search string, limit int) ([]*Task, error) {

	date, err := time.Parse("02.01.2006", search)
	if err == nil {
		// поиск по дате
		formatted := date.Format("20060102")
		return TasksByDate(formatted, limit)
	}

	pattern := "%" + search + "%"
	return TasksByPattern(pattern, limit)
}

func TasksByDate(formatted string, limit int) ([]*Task, error) {
	rows, err := DB.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date LIMIT ?", formatted, limit)
	if err != nil {
		return []*Task{}, err
	}

	defer rows.Close()
	tasks := []*Task{}

	for rows.Next() {
		task := &Task{}
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return []*Task{}, err
	}
	if tasks == nil {
		tasks = []*Task{}
	}

	return tasks, nil
}

func TasksByPattern(pattern string, limit int) ([]*Task, error) {
	rows, err := DB.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ? ", pattern, pattern, limit)
	if err != nil {
		return []*Task{}, err
	}

	defer rows.Close()
	tasks := []*Task{}

	for rows.Next() {
		task := &Task{}
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return []*Task{}, err
	}
	if tasks == nil {
		tasks = []*Task{}
	}

	return tasks, nil
}
