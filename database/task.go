package database

import (
	"database/sql"
	"errors"

	"main.go/model"
)

// AddTask добавляет задачу в базу данных, возвращает id задачи
func AddTask(db *sql.DB, task model.Task) (int64, error) {
	if task.Title == "" {
		return 0, errors.New("не указан заголовок задачи")
	}

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetTasks получает последние 50 задач из базы данных
func GetTasks(db *sql.DB) ([]model.Task, error) {
	var tasks []model.Task

	rows, err := db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT 50")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task model.Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// GetTask получает задачу по id
func GetTask(db *sql.DB, id string) (model.Task, error) {
	var task model.Task

	if err := db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
		if err == sql.ErrNoRows {
			return task, err
		}

	}
	return task, nil
}

// UpdateTask обновляет все поля задачи
func UpdateTask(db *sql.DB, task model.Task) error {

	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows != 1 {
		return errors.New("expected to affect 1 row")
	}
	return nil
}

// DeleteTask удаляет задачу из базы данных
func DeleteTask(db *sql.DB, id string) error {

	deleteQuery := `DELETE FROM scheduler WHERE id = ?`
	res, err := db.Exec(deleteQuery, id)
	if err != nil {

		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected != 1 {
		return errors.New("expected to affect 1 row")
	}

	return nil
}
