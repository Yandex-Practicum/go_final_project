package repository

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/AlexJudin/go_final_project/usecases/model"
)

var _ Task = (*TaskRepo)(nil)

type TaskRepo struct {
	Db *sqlx.DB
}

func NewNewRepository(db *sqlx.DB) *TaskRepo {
	return &TaskRepo{Db: db}
}

func (r *TaskRepo) CreateTask(task *model.Task) (int64, error) {
	res, err := r.Db.Exec(SQLCreateTask, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *TaskRepo) GetTasks() (model.TasksResp, error) {
	tasks := make([]model.Task, 0)

	res, err := r.Db.Query(SQLGetTasks, time.Now().Format("20060102"))
	if err != nil {
		return model.TasksResp{Tasks: tasks}, err
	}

	defer res.Close()

	var task model.Task

	for res.Next() {
		err = res.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return model.TasksResp{Tasks: tasks}, err
		}

		tasks = append(tasks, task)
	}

	return model.TasksResp{Tasks: tasks}, nil
}

func (r *TaskRepo) GetTasksBySearchString(searchString string) (model.TasksResp, error) {
	tasks := make([]model.Task, 0)

	res, err := r.Db.Query(SQLGetTasksBySearchString, searchString)
	if err != nil {
		return model.TasksResp{Tasks: tasks}, err
	}

	defer res.Close()

	var task model.Task

	for res.Next() {
		err = res.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return model.TasksResp{Tasks: tasks}, err
		}

		tasks = append(tasks, task)
	}

	return model.TasksResp{Tasks: tasks}, nil
}

func (r *TaskRepo) GetTasksByDate(searchString string) (model.TasksResp, error) {
	tasks := make([]model.Task, 0)

	res, err := r.Db.Query(SQLGetTasksByDate, searchString)
	if err != nil {
		return model.TasksResp{Tasks: tasks}, err
	}

	defer res.Close()

	var task model.Task

	for res.Next() {
		err = res.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return model.TasksResp{Tasks: tasks}, err
		}

		tasks = append(tasks, task)
	}

	return model.TasksResp{Tasks: tasks}, nil
}

func (r *TaskRepo) GetTaskById(id string) (*model.Task, error) {
	var task model.Task

	res, err := r.Db.Query(SQLGetTaskById, id)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	if res.Next() {
		err = res.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
	}

	if task.Id == "" {
		return nil, fmt.Errorf("task id %s not found", id)
	}

	return &task, nil
}

func (r *TaskRepo) UpdateTask(task *model.Task) error {
	_, err := r.Db.Exec(SQLUpdateTask, task.Id, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return err
	}

	return nil
}

func (r *TaskRepo) MakeTaskDone(id string, date string) error {
	_, err := r.Db.Exec(SQLMakeTaskDone, id, date)
	if err != nil {
		return err
	}

	return nil
}

func (r *TaskRepo) DeleteTask(id string) error {
	_, err := r.Db.Exec(SQLDeleteTask, id)
	if err != nil {
		return err
	}

	return nil
}
