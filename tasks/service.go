package tasks

import (
	"errors"
	"github.com/qavaleria/go_final_project/models"
	"time"
)

func ValidateTask(task *models.Task) error {
	if task.Title == "" {
		return errors.New("не указан заголовок задачи")
	}

	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format(FormatDate)
	} else {
		date, err := time.Parse(FormatDate, task.Date)
		if err != nil {
			return errors.New("Дата представлена в неправильном формате")
		}

		if date.Before(now) {
			if task.Repeat == "" {
				task.Date = now.Format(FormatDate)
			} else {
				nextDate, err := NextDate(now, task.Date, task.Repeat)
				if err != nil {
					return errors.New("Ошибка вычисления следующей даты")
				}
				task.Date = nextDate
			}
		}
	}

	return nil
}
