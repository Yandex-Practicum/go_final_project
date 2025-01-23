package rules

import (
	"errors"
	"time"
)

// handleYearRepeat обрабатывает ежегодное повторение и возвращает следующую дату в формате time.Time.
func handleYearRepeat(now time.Time, taskDate time.Time, rules []string) (string, error) {
	if len(rules) != 1 {
		return "", errors.New("неверный формат правидо повторения для года")
	}

	// Увеличиваем дату задачи до следующего года, если она не после текущей даты.
	if taskDate.After(now) {
		taskDate = taskDate.AddDate(1, 0, 0)
	} else {
		for !taskDate.After(now) {
			taskDate = taskDate.AddDate(1, 0, 0)
		}
	}

	// Возвращаем следующую дату
	return taskDate.Format("20060102"), nil
}
