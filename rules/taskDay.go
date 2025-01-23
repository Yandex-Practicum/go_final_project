package rules

import (
	"errors"
	"strconv"
	"time"
)

// handleDayRepeat обрабатывает правило повторения по дням.
func handleDayRepeat(now time.Time, taskDate time.Time, rules []string) (string, error) {

	if len(rules) != 2 {
		return "", errors.New("неверный формат правила повторения дня")
	}
	days, err := strconv.Atoi(rules[1])
	if err != nil || days <= 0 || days > 400 {
		return "", errors.New("некорректный интервал для дня")
	}

	// Если taskDate уже прошла, устанавливаем его на now
	if taskDate.After(now) {
		taskDate = taskDate.AddDate(0, 0, days)
	} else {
		for !taskDate.After(now) {
			taskDate = taskDate.AddDate(0, 0, days)
		}
	}

	return taskDate.Format("20060102"), nil
}
