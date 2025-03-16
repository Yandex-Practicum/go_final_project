package repository

import (
	"database/sql"
	"fmt"
	"test/internal/task"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{
		db: db,
	}
}

func (ts *TaskRepository) Insert(task *task.Task) (string, error) {
	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES ($1, $2, $3, $4)"
	res, err := ts.db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return "", err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return "", err
	}

	return fmt.Sprint(id), nil
}

func (ts *TaskRepository) GetAll() (*task.List, error) {
	query := "SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date"
	rows, err := ts.db.Query(query)
	if err != nil {
		return nil, err
	}

	return ts.prepareTaskList(rows)
}

func (ts *TaskRepository) GetByDate(date string) (*task.List, error) {
	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE `date` = ?"
	rows, err := ts.db.Query(query, date)
	if err != nil {
		return nil, err
	}

	return ts.prepareTaskList(rows)
}

func (ts *TaskRepository) GetByTitleOrComment(search string) (*task.List, error) {
	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE (title LIKE ? OR comment LIKE ?) ORDER BY date"
	rows, err := ts.db.Query(query, "%"+search+"%", "%"+search+"%")
	if err != nil {
		return nil, err
	}

	return ts.prepareTaskList(rows)
}

func (ts *TaskRepository) GetById(id int) (*task.Task, error) {
	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?"
	row := ts.db.QueryRow(query, id)
	var t task.Task
	err := row.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (ts *TaskRepository) prepareTaskList(rows *sql.Rows) (*task.List, error) {
	taskList := make([]*task.Task, 0)
	defer rows.Close()

	for rows.Next() {
		var taskStruct task.Task
		err := rows.Scan(&taskStruct.ID, &taskStruct.Date, &taskStruct.Title, &taskStruct.Comment, &taskStruct.Repeat)
		if err != nil {
			return nil, err
		}
		taskList = append(taskList, &taskStruct)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	return &task.List{Task: taskList}, nil
}

func (ts *TaskRepository) DeleteById(id int) error {
	query := "DELETE FROM scheduler WHERE id = $1"
	_, err := ts.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

func (ts *TaskRepository) UpdateById(t *task.Task) (*task.Task, error) {
	query := "UPDATE scheduler SET date = $1, title = $2, comment = $3, repeat = $4 WHERE id = $5"
	_, err := ts.db.Exec(query, t.Date, t.Title, t.Comment, t.Repeat, t.ID)
	if err != nil {
		return nil, err
	}

	return t, err
}

func (ts *TaskRepository) Done(task *task.Task) error {
	query := "UPDATE scheduler SET date = $1 WHERE id = $2"
	_, err := ts.db.Exec(query, task.Date, task.ID)
	if err != nil {
		return err
	}

	return nil
}
