package service

import (
	"strings"
	"time"

	"go_final_project/internal/models"
	"go_final_project/internal/utils"
)

type TaskRepository interface {
	CreateTask(task *models.Task) (int64, error)
	GetTaskByID(ID int64) (*models.Task, error)
	UpdateTask(task *models.Task) error
	UpdateTaskDate(taskID int64, date string) error
	DeleteTaskByID(ID int64) error
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

func (s *TaskService) GetTask(ID int64) (*models.Task, error) {
	return s.repository.GetTaskByID(ID)
}

func (s *TaskService) GetTasksWithFilter(filterType int, filterValue string) ([]*models.Task, error) {
	var tasks []*models.Task
	var err error
	switch filterType {
	case utils.FilterTypeDate:
		tasks, err = s.repository.GetAllTasksFilterByDate(filterValue)
	case utils.FilterTypeSearch:
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

	validTask.ID, err = s.repository.CreateTask(validTask)
	if err != nil {
		return 0, err
	}

	return validTask.ID, nil
}

func (s *TaskService) UpdateTask(task *models.Task) error {
	task, err := validateTask(task)
	if err != nil {
		return err
	}

	return s.repository.UpdateTask(task)
}

func (s *TaskService) DeleteTask(ID int64) error {
	return s.repository.DeleteTaskByID(ID)
}

func (s *TaskService) SetTaskDone(taskID int64) error {
	task, err := s.repository.GetTaskByID(taskID)
	if err != nil {
		return err
	}

	if len(task.Repeat) == 0 {
		return s.repository.DeleteTaskByID(task.ID)
	}

	task.Date, err = utils.NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		return err
	}

	return s.repository.UpdateTaskDate(task.ID, task.Date)
}

func validateTask(task *models.Task) (*models.Task, error) {
	if task.Title == "" {
		return nil, utils.ErrInvalidTaskTitle
	}

	now := time.Now()
	today := now.Format(utils.ParseDateFormat)
	if len(strings.TrimSpace(task.Date)) == 0 {
		task.Date = today
		return task, nil
	}

	taskDate, err := time.Parse(utils.ParseDateFormat, task.Date)
	if err != nil {
		return nil, utils.ErrInvalidTaskDate
	}

	if taskDate.Format(utils.ParseDateFormat) < today {
		if len(strings.TrimSpace(task.Repeat)) == 0 {
			task.Date = today
		} else {
			nextDate, err := utils.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return nil, utils.ErrInvalidTaskRepeat
			}
			task.Date = nextDate
		}
	}

	return task, nil
}
