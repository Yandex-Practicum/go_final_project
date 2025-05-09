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
