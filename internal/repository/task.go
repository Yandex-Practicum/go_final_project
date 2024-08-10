package repository

import (
	"database/sql"
	"errors"
	"go_final_project/internal/models"
	"go_final_project/internal/utils"
)

const getLimit = 50

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) CreateTask(task *models.Task) (int64, error) {
	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)"
	res, err := r.db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, errors.Join(utils.ErrDBInsert, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, errors.Join(utils.ErrGetTaskID, err)
	}
	return id, nil
}

func (r *TaskRepository) GetTaskByID(ID int64) (*models.Task, error) {
	query := `SELECT 
    			id,
    			date,
    			title,
    			comment,
    			repeat
			  FROM scheduler
			  WHERE id = ?`
	row := r.db.QueryRow(query, ID)
	var task models.Task
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.ErrTaskNotFound
		}
		return nil, errors.Join(utils.ErrTaskParse, err)
	}
	return &task, nil
}

func (r *TaskRepository) UpdateTask(task *models.Task) error {
	query := `UPDATE scheduler
    		  SET
    			date = ?,
    			title = ?,
    			comment = ?,
    			repeat = ?
			  WHERE id = ?`
	updateResult, err := r.db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return utils.ErrTaskNotFound
	}

	rowsAffected, err := updateResult.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return utils.ErrTaskNotFound
	}

	return nil
}

func (r *TaskRepository) UpdateTaskDate(taskID int64, date string) error {
	query := `UPDATE scheduler SET date = ? WHERE id = ?`
	_, err := r.db.Exec(query, date, taskID)
	return err
}

func (r *TaskRepository) DeleteTaskByID(ID int64) error {
	deleteQuery := `DELETE FROM scheduler WHERE id = ?`
	_, deleteErr := r.db.Exec(deleteQuery, ID)
	return deleteErr
}

func (r *TaskRepository) GetAllTasks() ([]*models.Task, error) {
	query := `SELECT 
    			id,
    			date,
    			title,
    			comment,
    			repeat
			  FROM scheduler
			  ORDER BY date
			  LIMIT ?`
	rows, selectErr := r.db.Query(query, getLimit)
	if selectErr != nil {
		return nil, selectErr
	}
	return r.parseTasks(rows)
}

func (r *TaskRepository) GetAllTasksFilterByDate(date string) ([]*models.Task, error) {
	query := `SELECT 
    			id,
    			date,
    			title,
    			comment,
    			repeat
			  FROM scheduler
			  WHERE date = ?
			  ORDER BY date
			  LIMIT ?`
	rows, selectErr := r.db.Query(query, date, getLimit)
	if selectErr != nil {
		return nil, selectErr
	}
	return r.parseTasks(rows)
}

func (r *TaskRepository) GetAllTasksFilterByTitleOrComment(search string) ([]*models.Task, error) {
	query := `SELECT 
    			id,
    			date,
    			title,
    			comment,
    			repeat
			  FROM scheduler
			  WHERE title LIKE ? OR comment LIKE ?
			  ORDER BY date
			  LIMIT ?`
	rows, selectErr := r.db.Query(query, search, search, getLimit)
	if selectErr != nil {
		return nil, selectErr
	}

	return r.parseTasks(rows)
}

func (r *TaskRepository) parseTasks(rows *sql.Rows) ([]*models.Task, error) {
	result := make([]*models.Task, 0)
	for rows.Next() {
		var selectTask models.Task
		err := rows.Scan(&selectTask.ID, &selectTask.Date, &selectTask.Title, &selectTask.Comment, &selectTask.Repeat)
		if err != nil {
			return result, errors.Join(utils.ErrTaskParse, err)
		}
		result = append(result, &selectTask)
	}
	return result, nil
}
