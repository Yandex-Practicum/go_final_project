package db

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

func GetTasks(s *sqlx.DB, search string) ([]Task, error) {

	const limit = 50

	var rows *sql.Rows

	var err error

	if search != "" {
		if date, err := time.Parse("02.01.2006", search); err == nil {
			// Формат даты в базе данных: 20060102
			dateStr := date.Format(DateFormat)
			rows, err = s.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE date =? LIMIT?", dateStr, limit)
		} else {
			rows, err = s.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE? OR comment LIKE? ORDER BY date LIMIT?", "%"+search+"%", "%"+search+"%", limit)
		}
	} else {
		rows, err = s.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT?", limit)
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
