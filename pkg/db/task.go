package db

import (
	"database/sql"
	"fmt"
)

// структура для Джейсона
type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// функция добавления задачи
func AddTask(task *Task) (int64, error) {
	var id int64
	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)"
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err == nil {
		id, err = res.LastInsertId()
	}

	return id, err
}

// функция отображения задач
func Tasks(limit int) ([]*Task, error) {
	rows, err := db.Query(`SELECT id, date, title, comment, repeat
	FROM scheduler
	ORDER BY date ASC
	LIMIT ?`, limit)
	if err != nil {
		return nil, fmt.Errorf("error while SELECT %w", err)
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, fmt.Errorf("error searching for nearby tasks: %w", err)
		}
		tasks = append(tasks, &t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error searching for nearby tasks: %w", err)
	}

	if tasks == nil {
		tasks = make([]*Task, 0)
	}
	return tasks, nil
}

// функция получения задачи
func GetTask(id string) (*Task, error) {
	if id == "" {
		return nil, fmt.Errorf("error scanning task by ID: %s", id)
	}

	t := &Task{}
	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id=?"
	err := db.QueryRow(query, id).Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task not found: %s", id)
		}
		return nil, fmt.Errorf("error returning rows: %w", err)
	}
	return t, nil
}

// функция редактирования задачи
func UpdateTask(t *Task) error {
	query := `
UPDATE scheduler
SET date=?,
    title=?,
    comment=?,
    repeat=?
    WHERE id=?`
	res, err := db.Exec(query, t.Date, t.Title, t.Comment, t.Repeat, t.ID)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if n == 0 {
		return fmt.Errorf("incorrect ID for task: %s", t.ID)
	}
	return nil

}

// функция удаления задачи
func DeleteTask(id string) error {
	if id == "" {
		return fmt.Errorf("ID not specified: %s", id)
	}
	res, err := db.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("error deleting task: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error counting deleted rows: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("task not found: %s", id)
	}
	return nil
}

// функция переноса даты
func UpdateDate(nextDate, id string) error {
	if id == "" {
		return fmt.Errorf("ID not specified")
	}
	res, err := db.Exec(`UPDATE scheduler SET date = ? WHERE id = ?`, nextDate, id)
	if err != nil {
		return fmt.Errorf("error updating date: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error counting updated rows: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}
