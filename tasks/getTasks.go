package tasks

import (
	"database/sql"
	"fmt"
	"time"
)

func GetTasks(dbPath, search string) ([]Task, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе данных: %w", err)
	}
	defer db.Close()

	const limit = 50

	var rows *sql.Rows

	if search != "" {
		if date, err := time.Parse("02.01.2006", search); err == nil {
			// Формат даты в базе данных: 20060102
			dateStr := date.Format("20060102")
			rows, err = db.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE date =? LIMIT?", dateStr, limit)
		} else {
			rows, err = db.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE? OR comment LIKE? ORDER BY date LIMIT?", "%"+search+"%", "%"+search+"%", limit)
		}
	} else {
		rows, err = db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT?", limit)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]Task, 0)
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
