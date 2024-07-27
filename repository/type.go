package repository

import "github.com/AlexJudin/go_final_project/usecases/model"

type Task interface {
	CreateTask(task *model.Task) (int64, error)
	GetTasks(searchString string) (model.TasksResp, error)
	GetTaskById(id string) (*model.Task, error)
	UpdateTask(task *model.Task) error
	MakeTaskDone(id string, date string) error
	DeleteTask(id string) error
}
