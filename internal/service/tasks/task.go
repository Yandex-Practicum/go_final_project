package tasks

import (
	"cactus3d/go_final_project/internal/models"
	"cactus3d/go_final_project/internal/utils"
	"fmt"
	"regexp"
	"time"
)

var (
	defaultLimit = 10
)

type Service struct {
	store StoreProvider
}

type StoreProvider interface {
	AddTask(task *models.Task) (int, error)
	GetTaskById(id string) (*models.Task, error)
	GetTasks(offset, limit int) ([]models.Task, error)
	GetTasksBySearch(search string, offset, limit int) ([]models.Task, error)
	GetTasksByDate(date string, offset, limit int) ([]models.Task, error)
	UpdateTask(*models.Task) (int64, error)
	DeleteTaskById(id string) (int64, error)
}

func New(store StoreProvider) *Service {
	return &Service{store: store}
}

func (s *Service) AddTask(date, title, comment, repeat string) (int, error) {
	_, err := time.Parse("20060102", date)
	if err != nil {
		return 0, err
	}
	n := time.Now().Format("20060102")

	if date < n {
		if repeat != "" {
			date, err = utils.NextDate(time.Now(), date, repeat)
			if err != nil {
				return 0, err
			}
		} else {
			date = time.Now().Format("20060102")
		}
	}

	task := &models.Task{
		Title:   title,
		Date:    date,
		Comment: comment,
		Repeat:  repeat,
	}

	return s.store.AddTask(task)
}

func (s *Service) GetTasks(search string, offset, limit int) ([]models.Task, error) {
	if limit == 0 {
		limit = defaultLimit
	}

	if search != "" {
		_, err := regexp.Compile(`[0-3][0-9]\.[0-1][0-9]\.20[0-9][0-9]`)
		if err != nil {
			return nil, err
		}
		matched, _ := regexp.MatchString(`[0-3][0-9]\.[0-1][0-9]\.20[0-9][0-9]`, search)
		if matched {
			date, err := time.Parse("02.01.2006", search)
			if err != nil {
				return nil, err
			}
			search = date.Format("20060102")
			return s.store.GetTasksByDate(search, offset, limit)
		}

		search = "%" + search + "%"
		return s.store.GetTasksBySearch(search, offset, limit)
	}

	return s.store.GetTasks(offset, limit)
}

func (s *Service) GetTask(id string) (*models.Task, error) {
	task, err := s.store.GetTaskById(id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, fmt.Errorf("not found")
	}

	return task, nil
}

func (s *Service) UpdateTask(task *models.Task) error {
	count, err := s.store.UpdateTask(task)
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("no rows with such id")
	}
	return nil
}

func (s *Service) DeleteTask(id string) error {
	count, err := s.store.DeleteTaskById(id)
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("nothing to delete")
	}

	return nil
}

func (s *Service) DoneTask(id string) error {
	task, err := s.GetTask(id)
	if err != nil {
		return err
	}

	if task.Repeat == "" {
		err = s.DeleteTask(task.Id)
		if err != nil {
			return err
		}
		return nil
	}

	date, err := utils.NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		return err
	}

	task.Date = date
	err = s.UpdateTask(task)
	if err != nil {
		return err
	}

	return nil
}
