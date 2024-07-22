package usecases

import (
	"time"

	"github.com/AlexJudin/go_final_project/usecases/model"
)

type Task interface {
	GetNextDate(now time.Time, date string, repeat string) (string, error)
	CreateTask(task *model.Task) (*model.TaskResp, error)
	GetTasks() (model.TasksResp, error)
	GetTaskById(id string) (*model.Task, error)
}
