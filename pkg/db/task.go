package db

import (
	"fmt"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat"`
}

type TaskResp struct {
	Tasks []*Task `json:"tasks"`
}

func AddTask(t Task) (int64, error) {
	var id int64

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`

	res, err := db.Exec(query, t.Date, t.Title, t.Comment, t.Repeat)
	if err == nil {
		id, err = res.LastInsertId()
	}

	return id, err
}

func GetTask(id string) (*Task, error) {
	query := `SELECT * FROM scheduler WHERE id = ?`
	row := db.QueryRow(query, id)

	var task Task

	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, fmt.Errorf("Could not read row %w\n", err)
	}

	return &task, nil
}

func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ?`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("Incorrect id for updating task")
	}

	return nil
}

func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = ?`
	_, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("Could not delete Task %w\n", err)
	}

	return nil
}

func UpdateDate(id, date string) error {
	query := `UPDATE scheduler SET date = ? WHERE id = ?`
	_, err := db.Exec(query, date, id)
	if err != nil {
		return fmt.Errorf("Could not update date %w\n", err)
	}

	return nil
}

func Tasks(limit int, search, tip string) ([]*Task, error) {
	var query string
	var args []interface{}

	switch tip {
	case "time":
		query = `SELECT * FROM scheduler WHERE date = ? ORDER BY date ASC LIMIT ?`
		args = []interface{}{search, limit}
	case "default":
		query = `SELECT * FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date ASC LIMIT ?`
		pattern := "%" + search + "%"
		args = []interface{}{pattern, pattern, limit}
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("Could not read rows from table %w\n", err)
	}
	defer rows.Close()

	tasks := []*Task{}

	for rows.Next() {
		var task Task

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, fmt.Errorf("Could not read row %w\n", err)
		}

		tasks = append(tasks, &task)
	}

	return tasks, nil
}
