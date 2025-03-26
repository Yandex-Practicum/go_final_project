package database

import (
	"log"

	"go_final_project/tasks"
)

func (d Database) GetTaskByID(id int) (tasks.Task, error) {
	var t tasks.Task
	log.Printf("%v", id)
	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?"
	err := d.db.QueryRow(query, id).Scan(&t.Id, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		return t, err
	}
	return t, nil
}
