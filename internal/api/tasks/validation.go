package tasks

import (
	"fmt"
	"time"

	"final_project/internal/api"
	"final_project/internal/db"
)

func validateAndFormatTask(now time.Time, task *db.Task) error {
	if task.Title == "" {
		return fmt.Errorf("не указан заголовок задачи")
	}

	if task.Date == "" {
		task.Date = now.Format(api.DateFormat)
	}

	taskDate, err := time.Parse(api.DateFormat, task.Date)
	if err != nil {
		return fmt.Errorf("некорректный формат даты")
	}

	if task.Repeat != "" {
		next, err := api.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
		if api.AfterNow(now, taskDate) {
			task.Date = next
		}
	} else {
		if api.AfterNow(now, taskDate) {
			task.Date = now.Format(api.DateFormat)
		}
	}

	return nil
}
