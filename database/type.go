package database

import "github.com/AlexJudin/go_final_project/usecases/model"

type Task interface {
	CreateTask(task *model.TaskReq) (int64, error)
}
