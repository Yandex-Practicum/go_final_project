package api

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// NextDate вычисляет следующую дату задачи на основе правила повторения
func NextDate(now time.Time, date string, repeat string) (string, error) {
	const layout = "20060102"
	taskDate, err := time.Parse(layout, date)
	if err != nil {
		return "", fmt.Errorf("неправильная дата: %v", err)
	}

	if repeat == "" {
		return "", errors.New("правило повторения не указано")
	}

	parts := strings.Fields(repeat)
	if len(parts) < 1 {
		return "", errors.New("неправильное правило повторения")
	}

	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", errors.New("неправильное правило для дня")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days < 1 || days > 400 {
			return "", errors.New("неправильное количество дней")
		}
		for {
			taskDate = taskDate.AddDate(0, 0, days)
			if taskDate.After(now) {
				return taskDate.Format(layout), nil
			}
		}
	case "y":
		taskDate = taskDate.AddDate(1, 0, 0)
		for !taskDate.After(now) {
			taskDate = taskDate.AddDate(1, 0, 0)
		}
		return taskDate.Format(layout), nil
	default:
		return "", errors.New("неподдерживаемое правило повторения")
	}
}
