package db

import (
	"database/sql"
	"fmt"
	"pwd/internal/controller"

	_ "modernc.org/sqlite"
)

func New() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		return nil, err
	}

	return db, nil
}

func GetTasks(db *sql.DB) ([]controller.Task, error) {
	rows, err := db.Query(
		"SELECT id, date, title, comment, repeat " +
			"FROM scheduler ORDER BY date DESC LIMIT 10 OFFSET 0",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]controller.Task, 0)
	for rows.Next() {
		var task controller.Task
		if err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {

			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func AddTask(db *sql.DB, task controller.Task) (int, error) {
	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)"
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		fmt.Println(err)
	}

	return int(id), nil
}
