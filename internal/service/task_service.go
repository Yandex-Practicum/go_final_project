package service

import (
	"errors"
	"go_final_project/internal/model"
	"go_final_project/internal/repository"
	"strconv"
	"strings"
	"time"
)

type TaskService interface {
	AddTask(task model.Task) (int64, error)
	GetTasksByDate(date string, limit int) ([]model.Task, error)
	GetTasksBySearch(search string, limit int) ([]model.Task, error)
	GetAllTasks(limit int) ([]model.Task, error)
	GetTaskByID(id string) (model.Task, error)
	UpdateTask(task model.Task) error
	DeleteTask(id int) error
	TaskDone(id int) error
}

type taskService struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) TaskService {
	return &taskService{repo: repo}
}

func (s *taskService) AddTask(task model.Task) (int64, error) {
	if task.Title == "" {
		return 0, errors.New("не указан заголовок задачи")
	}

	now := time.Now()
	todayStr := now.Format("20060102")

	if task.Date == "" {
		task.Date = todayStr
	}

	taskDate, err := time.Parse("20060102", task.Date)
	if err != nil {
		return 0, errors.New("неправильный формат даты")
	}

	if taskDate.Before(now) {
		task.Date = todayStr
	}

	if task.Repeat != "" {
		_, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return 0, errors.New("неправильное правило повторения")
		}
	}

	return s.repo.AddTask(task)
}

func (s *taskService) GetTasksByDate(date string, limit int) ([]model.Task, error) {
	return s.repo.GetTasksByDate(date, limit)
}

func (s *taskService) GetTasksBySearch(search string, limit int) ([]model.Task, error) {
	return s.repo.GetTasksBySearch(search, limit)
}

func (s *taskService) GetAllTasks(limit int) ([]model.Task, error) {
	return s.repo.GetAllTasks(limit)
}

func (s *taskService) GetTaskByID(id string) (model.Task, error) {
	return s.repo.GetTaskByID(id)
}

func (s *taskService) UpdateTask(task model.Task) error {
	if task.ID == "" {
		return errors.New("не указан идентификатор задачи")
	}

	_, err := strconv.Atoi(task.ID)
	if err != nil {
		return errors.New("неверный формат идентификатора задачи")
	}

	_, err = s.repo.GetTaskByID(task.ID)
	if err != nil {
		return errors.New("задача не найдена")
	}

	if task.Title == "" {
		return errors.New("не указан заголовок задачи")
	}

	now := time.Now()
	todayStr := now.Format("20060102")

	if task.Date == "" {
		task.Date = todayStr
	}

	taskDate, err := time.Parse("20060102", task.Date)
	if err != nil {
		return errors.New("неправильный формат даты")
	}

	if taskDate.Before(now) {
		task.Date = todayStr
	}

	if task.Repeat != "" {
		_, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return errors.New("неправильное правило повторения")
		}
	}

	return s.repo.UpdateTask(task)
}

func (s *taskService) DeleteTask(id int) error {
	return s.repo.DeleteTask(id)
}

func (s *taskService) TaskDone(id int) error {
	task, err := s.repo.GetTaskByID(strconv.Itoa(id))
	if err != nil {
		return err
	}

	if task.Repeat != "" {
		now := time.Now()
		nextDate, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
		task.Date = nextDate
		return s.repo.UpdateTask(task)
	} else {
		return s.repo.DeleteTask(id)
	}
}

func NextDate(now time.Time, date string, repeat string) (string, error) {
	d, err := time.Parse("20060102", date)
	if err != nil {
		return "", errors.New("неправильный формат даты")
	}

	if repeat == "" {
		return "", errors.New("пустое правило повторения")
	}

	parts := strings.Fields(repeat)
	if len(parts) < 1 {
		return "", errors.New("неверный формат повтора")
	}

	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", errors.New("неверный d формат")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days < 1 || days > 400 {
			return "", errors.New("неправильное значение дней")
		}

		if d.Before(now) || d.Equal(now) {
			for d.Before(now) || d.Equal(now) {
				d = d.AddDate(0, 0, days)
			}
		} else {
			// Если дата d еще не наступила, просто добавляем дни
			d = d.AddDate(0, 0, days)
		}

		// }
		return d.Format("20060102"), nil

	case "y":
		nextDate := d.AddDate(1, 0, 0)
		for !nextDate.After(now) {
			nextDate = nextDate.AddDate(1, 0, 0)
		}
		return nextDate.Format("20060102"), nil

	case "w":
		if len(parts) != 2 {
			return "", errors.New("неверный w формат")
		}
		daysStr := strings.Split(parts[1], ",")
		days := make([]int, len(daysStr))
		for i, dayStr := range daysStr {
			day, err := strconv.Atoi(dayStr)
			if err != nil || day < 1 || day > 7 {
				return "", errors.New("неправильный день недели")
			}
			days[i] = day
		}

		// Начинаем с следующего дня
		nextDate := d
		if nextDate.Before(now) {
			nextDate = now
		}
		for {
			nextDate = nextDate.AddDate(0, 0, 1)
			weekday := int(nextDate.Weekday())
			if weekday == 0 {
				weekday = 7
			}
			for _, day := range days {
				if weekday == day {
					return nextDate.Format("20060102"), nil
				}
			}
		}

	case "m":
		if len(parts) < 2 {
			return "", errors.New("неверный m формат")
		}

		daysStr := strings.Split(parts[1], ",")
		days := make([]int, len(daysStr))
		for i, dayStr := range daysStr {
			day, err := strconv.Atoi(dayStr)
			if err != nil || (day < -2 || day > 31) || day == 0 {
				return "", errors.New("неправильное значение дней")
			}
			days[i] = day
		}

		var months []int
		if len(parts) == 3 {
			monthsStr := strings.Split(parts[2], ",")
			months = make([]int, len(monthsStr))
			for i, monthStr := range monthsStr {
				month, err := strconv.Atoi(monthStr)
				if err != nil || month < 1 || month > 12 {
					return "", errors.New("неверное значение месяца")
				}
				months[i] = month
			}
		}

		nextDate := d
		for {
			nextDate = nextDate.AddDate(0, 0, 1)
			for _, day := range days {
				var dayOfMonth int
				switch day {
				case -1:
					dayOfMonth = getLastDayOfMonth(nextDate)
				case -2:
					dayOfMonth = getLastDayOfMonth(nextDate) - 1
				default:
					dayOfMonth = day
				}
				if nextDate.Day() == dayOfMonth {
					if len(months) == 0 || contains(months, int(nextDate.Month())) {
						if nextDate.After(now) {
							return nextDate.Format("20060102"), nil
						}
					}
				}
			}
		}

	default:
		return "", errors.New("неверный тип повтора")
	}

}

func getLastDayOfMonth(t time.Time) int {
	return time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func contains(slice []int, val int) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
