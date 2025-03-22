package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("Параметр repeat не может быть пустым")
	}

	taskDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", errors.New("Неверный формат даты")
	}

	parts := strings.Fields(repeat)

	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", errors.New("Неверный формат repeat: 'd' должен содержать одно число")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", errors.New("Неверное значение для количества дней")
		}
		if days < 1 || days > 400 {
			return "", errors.New("Количество дней должно быть от 1 до 400 для 'd'")
		}
		taskDate = taskDate.AddDate(0, 0, days)

	case "y":
		if len(parts) > 1 {
			return "", errors.New("Неверный формат repeat: 'y' не должен содержать аргументы")
		}
		taskDate = taskDate.AddDate(1, 0, 0)

	default:
		return "", fmt.Errorf("Неверный формат repeat: недопустимый символ '%s'", parts[0])
	}

	for !taskDate.After(now) {
		switch parts[0] {
		case "d":
			days, _ := strconv.Atoi(parts[1])
			taskDate = taskDate.AddDate(0, 0, days)
		case "y":
			taskDate = taskDate.AddDate(1, 0, 0)
		}
	}

	return taskDate.Format("20060102"), nil
}
