package repository

import (
	"database/sql"
	"final_project/internal/common"
)

type Tasks struct {
	Tasks []common.Task `json:"tasks"`
}

func (rep *Repository) GetAllTasks() (tasks Tasks, err error) {
	t := Tasks{}
	query := "SELECT * FROM scheduler ORDER BY date LIMIT :limit "
	rows, err := rep.db.Query(query, sql.Named("limit", common.Limit))
	if err != nil {
		return t, err
	}
	defer rows.Close()
	for rows.Next() {
		task := common.Task{}
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return t, err
		}
		tasks.Tasks = append(tasks.Tasks, task)
	}
	return tasks, nil
}

func (rep *Repository) GetTasksByDate(date string) (tasks Tasks, err error) {
	t := Tasks{}
	query := "SELECT * FROM scheduler WHERE date=:param LIMIT :limit"
	rows, err := rep.db.Query(query, sql.Named("param", date), sql.Named("limit", common.Limit))
	if err != nil {
		return t, err
	}
	defer rows.Close()
	for rows.Next() {
		task := common.Task{}
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return t, err
		}
		tasks.Tasks = append(tasks.Tasks, task)
	}
	return tasks, nil
}

func (rep *Repository) GetTasksByParam(param string) (tasks Tasks, err error) {
	t := Tasks{}
	query := "SELECT * FROM scheduler WHERE title LIKE :param OR comment LIKE :param ORDER BY date LIMIT :limit"
	rows, err := rep.db.Query(query, sql.Named("param", "%"+param+"%"), sql.Named("limit", common.Limit))
	if err != nil {
		return t, err
	}
	defer rows.Close()
	for rows.Next() {
		task := common.Task{}
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return t, err
		}
		tasks.Tasks = append(tasks.Tasks, task)
	}
	return tasks, nil
}
