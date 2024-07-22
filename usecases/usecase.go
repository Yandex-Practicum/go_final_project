package usecases

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/AlexJudin/go_final_project/repository"
	"github.com/AlexJudin/go_final_project/usecases/model"
)

var _ Task = (*TaskUsecase)(nil)

const (
	year = 1
)

type TaskUsecase struct {
	DB repository.Task
}

func NewTaskUsecase(db repository.Task) *TaskUsecase {
	return &TaskUsecase{DB: db}
}

func (t *TaskUsecase) GetNextDate(now time.Time, date string, repeat string) (string, error) {
	dateTask, err := time.Parse("20060102", date)
	if err != nil {
		return "", err
	}

	dateTaskNow := time.Now().Format("20060102")
	if repeat == "" {
		return dateTaskNow, nil
	}

	repeatString := strings.Split(repeat, " ")

	switch strings.ToLower(repeatString[0]) {
	case "d":
		days, err := parseValue(repeatString[1])
		if err != nil {
			return "", err
		}
		if days == 1 {
			return dateTaskNow, nil
		}
		dateTask = addDateTask(now, dateTask, 0, 0, days)
	case "y":
		dateTask = addDateTask(now, dateTask, year, 0, 0)
	//case "w":
	//case "m":
	default:
		return "", fmt.Errorf("invalid character")
	}

	return dateTask.Format("20060102"), nil
}

func (t *TaskUsecase) CreateTask(task *model.Task) (*model.TaskResp, error) {
	taskId, err := t.DB.CreateTask(task)
	if err != nil {
		return nil, err
	}

	taskResp := model.NewTaskResp(taskId)

	return taskResp, nil
}

func (t *TaskUsecase) GetTasks() (model.TasksResp, error) {
	return t.DB.GetTasks()
}

func (t *TaskUsecase) GetTaskById(id string) (*model.Task, error) {
	return t.DB.GetTaskById(id)
}

func parseValue(num string) (int, error) {
	days, err := strconv.Atoi(num)
	if err != nil {
		return 0, err
	}

	if days >= 400 || days < 0 {
		return 0, fmt.Errorf("invalid value %d", days)
	}

	return days, nil
}

func addDateTask(now time.Time, dateTask time.Time, year int, month int, day int) time.Time {
	dateTask = dateTask.AddDate(year, month, day)

	for dateTask.Before(now) {
		dateTask = dateTask.AddDate(year, month, day)
	}

	return dateTask
}
