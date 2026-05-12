package service

import (
	"errors"
	"go_final_project/internal/models"
	"go_final_project/pkg/api"
	"log"
	"strconv"
	"time"
)

var (
	ErrTaskIDEmpty     = errors.New("task id empty")
	ErrTaskIDInvalid   = errors.New("task id invalid")
	ErrTaskTitleEmpty  = errors.New("task title empty")
	ErrTaskDateInvalid = errors.New("task date invalid")
)

type ITaskStore interface {
	List(int) ([]models.Task, error)
	Get(string) (*models.Task, error)
	Add(*models.Task) (string, error)
	Update(*models.Task) error
	UpdateDate(*models.Task) error
	Delete(string) error
}

type TaskService struct {
	store ITaskStore
	logs  *log.Logger
}

func NewTaskService(store ITaskStore, logs *log.Logger) *TaskService {
	return &TaskService{
		store: store,
		logs:  logs,
	}
}

func (s *TaskService) Add(task *models.Task) (string, error) {
	if task.Title == "" {
		s.logs.Println(ErrTaskTitleEmpty)

		return "", ErrTaskTitleEmpty
	}

	now := time.Now()

	if task.Date == "" {
		task.Date = api.FormatDate(now)
	}

	taskDate, err := api.ParseDate(task.Date)

	if err != nil {
		s.logs.Println(err)

		return "", ErrTaskDateInvalid
	}

	if task.Repeat != "" {
		date, err := api.NextDate(now, task.Date, task.Repeat)

		if err != nil {
			s.logs.Println(err)

			return "", err
		}

		task.Date = date
	}

	if taskDate.Before(now) {
		task.Date = api.FormatDate(now)
	}

	return s.store.Add(task)
}

func (s *TaskService) List(limit int) ([]models.Task, error) {
	return s.store.List(limit)
}

func (s *TaskService) Get(id string) (*models.Task, error) {
	err := s.validateId(id)
	if err != nil {
		return nil, err
	}

	return s.store.Get(id)
}

func (s *TaskService) Edit(task *models.Task) error {
	err := s.validateId(task.ID)
	if err != nil {
		return err
	}

	if task.Title == "" {
		s.logs.Println(ErrTaskTitleEmpty)

		return ErrTaskTitleEmpty
	}

	_, err = api.ParseDate(task.Date)
	if err != nil {
		s.logs.Println(err)

		return ErrTaskDateInvalid
	}

	if task.Repeat != "" {
		_, err := api.NextDate(time.Now(), task.Date, task.Repeat)

		if err != nil {
			s.logs.Println(err)

			return err
		}
	}

	return s.store.Update(task)
}

func (s *TaskService) Delete(id string) error {
	err := s.validateId(id)
	if err != nil {
		return err
	}

	return s.store.Delete(id)
}

func (s *TaskService) Done(id string) error {
	task, err := s.Get(id)
	if err != nil {
		return err
	}

	if task.Repeat == "" {
		return s.Delete(task.ID)
	}

	task.Date, err = s.NextDate(task.Date, task.Date, task.Repeat)
	if err != nil {
		return err
	}

	return s.store.UpdateDate(task)
}

func (s *TaskService) NextDate(nowstr, date, repeat string) (string, error) {
	now, err := api.ParseDate(nowstr)
	if err != nil {
		s.logs.Println(err)
		return "", err
	}

	date, err = api.NextDate(now, date, repeat)
	if err != nil {
		s.logs.Println(err)
		return "", err
	}

	return date, nil
}

func (s *TaskService) validateId(id string) error {
	if id == "" {
		s.logs.Println(ErrTaskIDEmpty)

		return ErrTaskIDEmpty
	}

	_, err := strconv.Atoi(id)
	if err != nil {
		s.logs.Println(err)

		return ErrTaskIDInvalid
	}

	return nil
}
