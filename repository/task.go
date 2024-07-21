package repository

import (
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

func (r *TaskRepo) CreateTask(task *model.TaskReq) (int64, error) {
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
	tasks := make([]model.TaskReq, 0)

	res, err := r.Db.Query(SQLGetTasks, time.Now().Format("20060102"))
	if err != nil {
		return model.TasksResp{Tasks: tasks}, err
	}
	defer res.Close()

	var task model.TaskReq

	for res.Next() {
		err = res.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return model.TasksResp{Tasks: tasks}, err
		}

		tasks = append(tasks, task)
	}

	return model.TasksResp{Tasks: tasks}, nil
}
