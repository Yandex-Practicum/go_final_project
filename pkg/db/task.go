package db

import "fmt"

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

	result, err := DB.Exec(query,
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

func Tasks(limit int) ([]*Task, error) {
	query := `
		SELECT id, date, title, comment, repeat
		FROM scheduler
		ORDER BY date
		LIMIT ?
	`

	rows, err := DB.Query(query, limit)
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

func GetTask(id string) (*Task, error) {
	var task Task

	query := `
        SELECT id, date, title, comment, repeat
        FROM scheduler
        WHERE id = ?
    `

	err := DB.QueryRow(query, id).Scan(
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

	result, err := DB.Exec(
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
		return fmt.Errorf("Задача не найдена")
	}

	return nil
}

func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = ?`

	result, err := DB.Exec(query, id)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("Задача не найдена")
	}

	return nil
}

func UpdateDate(next string, id string) error {
	query := `UPDATE scheduler SET date = ? WHERE id = ?`

	result, err := DB.Exec(query, next, id)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("Задача не найдена")
	}

	return nil
}
