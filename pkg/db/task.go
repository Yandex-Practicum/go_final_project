package db

import (
	"fmt"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	query := `
		INSERT INTO scheduler (date, title, comment, repeat)
		VALUES (?, ?, ?, ?)
	`

	result, err := database.Exec(
		query,
		task.Date,
		task.Title,
		task.Comment,
		task.Repeat,
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

func GetTask(id string) (*Task, error) {
	query := `
		SELECT id, date, title, comment, repeat
		FROM scheduler
		WHERE id = ?
	`

	var task Task

	err := database.QueryRow(query, id).Scan(
		&task.ID,
		&task.Date,
		&task.Title,
		&task.Comment,
		&task.Repeat,
	)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func UpdateTask(task *Task) error {
	query := `
		UPDATE scheduler
		SET date = ?, title = ?, comment = ?, repeat = ?
		WHERE id = ?
	`

	result, err := database.Exec(
		query,
		task.Date,
		task.Title,
		task.Comment,
		task.Repeat,
		task.ID,
	)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}

func Tasks(limit int, search string) ([]*Task, error) {
	var (
		query  string
		params []any
	)

	if search == "" {
		query = `
			SELECT id, date, title, comment, repeat
			FROM scheduler
			ORDER BY date ASC
			LIMIT ?
		`

		params = []any{limit}

	} else if date, err := time.Parse("02.01.2006", search); err == nil {
		query = `
			SELECT id, date, title, comment, repeat
			FROM scheduler
			WHERE date = ?
			ORDER BY date ASC
			LIMIT ?
		`

		params = []any{
			date.Format("20060102"),
			limit,
		}

	} else {
		query = `
			SELECT id, date, title, comment, repeat
			FROM scheduler
			WHERE title LIKE ? OR comment LIKE ?
			ORDER BY date ASC
			LIMIT ?
		`

		searchPattern := "%" + search + "%"

		params = []any{
			searchPattern,
			searchPattern,
			limit,
		}
	}

	rows, err := database.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]*Task, 0)

	for rows.Next() {
		var task Task

		err := rows.Scan(
			&task.ID,
			&task.Date,
			&task.Title,
			&task.Comment,
			&task.Repeat,
		)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, &task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func DeleteTask(id string) error {
	query := `
		DELETE FROM scheduler
		WHERE id = ?
	`

	result, err := database.Exec(query, id)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}

func UpdateDate(next string, id string) error {
	query := `
		UPDATE scheduler
		SET date = ?
		WHERE id = ?
	`

	result, err := database.Exec(query, next, id)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}
