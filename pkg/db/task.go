package db

import (
	"database/sql"
	"errors"
	"fmt"
)

// Task представляет задачу
type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// AddTask добавляет задачу в базу данных
func AddTask(task *Task) (int64, error) {

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)`
	res, err := db.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	)

	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// GetTask возвращает задачу по ID
func GetTask(id string) (*Task, error) {

	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id`

	t := &Task{}

	if err := db.QueryRow(query, sql.Named("id", id)).Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("incorrect id for getting task")
		}
		return nil, err
	}
	return t, nil
}

// UpdateTask обновляет задачу
func UpdateTask(task *Task) error {

	query := `UPDATE scheduler SET date=:date, title=:title, comment=:comment, repeat=:repeat WHERE id=:id`
	res, err := db.Exec(query,
		sql.Named("id", task.ID),
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	)

	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("incorrect id for updating task")
	}
	return nil
}

// DeleteTask удаляет задачу
func DeleteTask(id string) error {

	query := `DELETE FROM scheduler WHERE id = :id`
	res, err := db.Exec(query, sql.Named("id", id))
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("incorrect id for deleting task")
	}
	return nil
}

// UpdateTaskDate обновляет дату задачи
func UpdateDate(id string, date string) error {
	if db == nil {
		return errors.New("database not initialized")
	}

	query := `UPDATE scheduler SET date = :date WHERE id = :id`
	res, err := db.Exec(query,
		sql.Named("date", date),
		sql.Named("id", id))

	if err != nil {
		return fmt.Errorf("update date failed: %w", err)
	}

	// Проверяем, что задача была удалена
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected error: %w", err)
	}
	if count == 0 {
		return errors.New("error in deliting task")
	}

	return nil
}
