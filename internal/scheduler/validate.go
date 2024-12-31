package scheduler

import (
	"errors"
	"go_final_project/internal/models"
	"time"
)

func ValidateAndProcessTask(task *models.Task, now time.Time) error {
	if task.Title == "" {
		return errors.New("загаловок не указан")
	}

	if task.Date == "" {
		task.Date = now.Format("20060102")
	}

	parsedDate, err := time.Parse("20060102", task.Date)
	if err != nil {
		return errors.New("неправильный формат date")
	}

	if startOfDay(parsedDate).Before(startOfDay(now)) {
		if task.Repeat == "" {
			task.Date = now.Format("20060102")
		} else {
			nextDate, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return errors.New("повторяющееся правило имеет неверный формат")
			}
			task.Date = nextDate
		}
	}

	return nil
}
