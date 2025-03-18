package services

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/wissio/go_final_project/internal/models"
)

const DateLayout = "20060102"

func ValidateTask(log *slog.Logger, task *models.Task) error {
	if task.Title == "" {
		return fmt.Errorf("empty title")
	}

	if task.Date == "" {
		task.Date = time.Now().Format(DateLayout)
	}

	validDate, err := time.Parse(DateLayout, task.Date)
	if err != nil {
		return fmt.Errorf("invalid date format")
	}

	if task.Repeat != "" && !isValidRepeatRule(task.Repeat) {
		return fmt.Errorf("invalid repeat rule")
	}

	if task.Repeat != "" && !isValidRepeatValue(task.Repeat) {
		return fmt.Errorf("invalid repeat rule")
	}

	currentDate := time.Now().Truncate(24 * time.Hour)
	validDate = validDate.Truncate(24 * time.Hour)

	if validDate.Before(currentDate) && !isValidDate(task.Repeat, currentDate, validDate) {
		if task.Repeat == "" {
			task.Date = time.Now().Format(DateLayout)
		}

		if task.Repeat != "" {
			task.Date, err = NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				return fmt.Errorf("error calculating next date: %v", err)
			}
		}
	}

	return nil
}

func isValidRepeatRule(repeat string) bool {
	return repeat[0] == 'd' || repeat[0] == 'w' || repeat[0] == 'm' || repeat[0] == 'y'
}

func isValidRepeatValue(repeat string) bool {
	if repeat[0] == 'd' || repeat[0] == 'w' || repeat[0] == 'm' {
		return len(repeat) >= 3
	}
	return true
}

func isValidDate(repeat string, currentDate time.Time, validDate time.Time) bool {
	if repeat == "" {
		return true
	}

	if validDate.After(currentDate) {
		return true
	}

	return false
}
