package repository

import (
	"database/sql"
	"final_project/internal/common"
)

func (rep *Repository) GetTaskByID(id string) (task common.Task, err error) {
	query := "SELECT * FROM scheduler WHERE id=:id"
	row := rep.db.QueryRow(query, sql.Named("id", id))

	err = row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return task, err
	}
	return task, nil
}
