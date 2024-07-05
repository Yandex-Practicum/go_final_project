package tasks

import (
	"database/sql"
)

func GetTaskByID(id string) (Task, error) {
	db, err := sql.Open("sqlite3", "./scheduler.db")
	if err != nil {
		return Task{}, err
	}
	defer db.Close()

	var task Task
	err = db.QueryRow("SELECT * FROM scheduler WHERE id = ?", id).Scan(
		&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat,
	)
	return task, err
}
