package database

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/qavaleria/go_final_project/models"
	"github.com/qavaleria/go_final_project/tasks"
	"time"
)

func AddTask(db *sqlx.DB, task models.Task) (int64, error) {
	if task.Title == "" {
		return 0, errors.New("Не указан заголовок задачи")
	}

	if task.Date == "" {
		task.Date = time.Now().Format(tasks.FormatDate)
	} else {
		_, err := time.Parse(tasks.FormatDate, task.Date)
		if err != nil {
			return 0, errors.New("Дата представлена в неправильном формате")
		}
	}

	now := time.Now()
	if task.Date < now.Format(tasks.FormatDate) {
		if task.Repeat == "" {
			task.Date = now.Format(tasks.FormatDate)
		} else {
			nextDate, err := tasks.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return 0, err
			}
			task.Date = nextDate
		}
	}

	result, err := db.Exec(
		`INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`,
		task.Date, task.Title, task.Comment, task.Repeat,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}
