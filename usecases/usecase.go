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

	repeatString := strings.Split(repeat, " ")

	switch strings.ToLower(repeatString[0]) {
	case "d":
		if len(repeatString) < 2 {
			return "", fmt.Errorf("repeat should be at least two characters for days")
		}

		days, err := parseValue(repeatString[1])
		if err != nil {
			return "", err
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

func (t *TaskUsecase) CreateTask(task *model.Task, pastDay bool) (*model.TaskResp, error) {
	if pastDay {
		nextDate, err := t.GetNextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			return nil, err
		}

		task.Date = nextDate
	}

	taskId, err := t.DB.CreateTask(task)
	if err != nil {
		return nil, err
	}

	taskResp := model.NewTaskResp(taskId)

	return taskResp, nil
}

func (t *TaskUsecase) GetTasks(searchString string) (model.TasksResp, error) {
	date, err := time.Parse("02.01.2006", searchString)
	if err == nil {
		return t.DB.GetTasksByDate(date)
	}

	if searchString != "" {
		return t.DB.GetTasksBySearchString(searchString)
	}

	return t.DB.GetTasks()
}

func (t *TaskUsecase) GetTaskById(id string) (*model.Task, error) {
	return t.DB.GetTaskById(id)
}

func (t *TaskUsecase) UpdateTask(task *model.Task, pastDay bool) error {
	_, err := t.GetTaskById(task.Id)
	if err != nil {
		return err
	}

	if pastDay {
		nextDate, err := t.GetNextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			return err
		}

		task.Date = nextDate
	}

	return t.DB.UpdateTask(task)
}

func (t *TaskUsecase) MakeTaskDone(id string) error {
	task, err := t.DB.GetTaskById(id)
	if err != nil {
		return err
	}

	if task.Repeat == "" {
		return t.DB.DeleteTask(id)
	}

	nextDate, err := t.GetNextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		return err
	}

	return t.DB.MakeTaskDone(id, nextDate)
}

func (t *TaskUsecase) DeleteTask(id string) error {
	_, err := t.DB.GetTaskById(id)
	if err != nil {
		return err
	}

	return t.DB.DeleteTask(id)
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
