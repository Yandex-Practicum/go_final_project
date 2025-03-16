package service

import (
	"strings"
	"time"

	"go_final_project/internal/constants"
	"go_final_project/internal/models"
	"go_final_project/internal/next_date"
)

type TaskRepository interface {
	CreateTask(task *models.Task) (int64, error)
	Get(id int64) (*models.Task, error)
	UpdateTask(task *models.Task) error
	UpdateTaskDate(taskId int64, date string) error
	DeleteTaskById(Id int64) error
	GetAllTasks() ([]*models.Task, error)
	GetAllTasksFilterByDate(date string) ([]*models.Task, error)
	GetAllTasksFilterByTitleOrComment(search string) ([]*models.Task, error)
}

type TaskService struct {
	repository TaskRepository
}

func NewTaskService(repository TaskRepository) *TaskService {
	return &TaskService{repository: repository}
}

func (s *TaskService) GetTask(id int64) (*models.Task, error) {
	return s.repository.Get(id)
}

func (s *TaskService) GetTasksWithFilter(filterType int, filterValue string) ([]*models.Task, error) {
	var tasks []*models.Task
	var err error
	switch filterType {
	case constants.FilterTypeDate:
		tasks, err = s.repository.GetAllTasksFilterByDate(filterValue)
	case constants.FilterTypeSearch:
		tasks, err = s.repository.GetAllTasksFilterByTitleOrComment(filterValue)
	default:
		tasks, err = s.repository.GetAllTasks()
	}
	return tasks, err
}

func (s *TaskService) CreateTask(task *models.Task) (int64, error) {
	validTask, err := validateTask(task)
	if err != nil {
		return 0, err
	}

	validTask.Id, err = s.repository.CreateTask(validTask)
	if err != nil {
		return 0, err
	}

	return validTask.Id, nil
}

func (s *TaskService) UpdateTask(task *models.Task) error {
	task, err := validateTask(task)
	if err != nil {
		return err
	}

	return s.repository.UpdateTask(task)
}

func (s *TaskService) DeleteTask(id int64) error {
	return s.repository.DeleteTaskById(id)
}

func (s *TaskService) SetTaskDone(id int64) error {
	task, err := s.repository.Get(id)
	if err != nil {
		return err
	}

	if len(task.Repeat) == 0 {
		return s.repository.DeleteTaskById(task.Id)
	}

	task.Date, err = next_date.NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		return err
	}

	return s.repository.UpdateTaskDate(task.Id, task.Date)
}

func validateTask(task *models.Task) (*models.Task, error) {
	if task.Title == "" {
		return nil, constants.ErrInvalidTaskTitle
	}

	now := time.Now()
	today := now.Format(constants.ParseDateFormat)
	if len(strings.TrimSpace(task.Date)) == 0 {
		task.Date = today
		return task, nil
	}

	taskDate, err := time.Parse(constants.ParseDateFormat, task.Date)
	if err != nil {
		return nil, constants.ErrInvalidTaskDate
	}

	if taskDate.Format(constants.ParseDateFormat) < today {
		if len(strings.TrimSpace(task.Repeat)) == 0 {
			task.Date = today
		} else {
			nextDate, err := next_date.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return nil, constants.ErrInvalidTaskRepeat
			}
			task.Date = nextDate
		}
	}

	return task, nil
}
