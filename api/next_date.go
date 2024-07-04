package api

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"main.go/model"
)

// NextDate считает следующую дату
func NextDate(now time.Time, date string, repeat string, layout string) (string, error) {
	parseDate, err := time.Parse(layout, date)
	if err != nil {
		return date, err
	}
	if repeat == "" {
		return date, errors.New("Пустая строка в repeat")
	}
	splitRepeat := strings.Split(repeat, " ")
	repeatDate := splitRepeat[0]

	switch repeatDate {
	case "d":
		if len(splitRepeat) < 2 {
			return date, errors.New("Не указан интервал в днях")
		}
		repeatNumber, err := strconv.Atoi(splitRepeat[1])
		if err != nil {
			return date, err
		}
		if repeatNumber > 400 {
			return date, errors.New("Превышен максимально допустимый интервал")
		}

		parseDate = parseDate.AddDate(0, 0, repeatNumber)

		for now.After(parseDate) {
			parseDate = parseDate.AddDate(0, 0, repeatNumber)
		}

	case "y":
		if len(splitRepeat) != 1 {
			return date, errors.New("Лишние данные")
		}
		parseDate = parseDate.AddDate(1, 0, 0)

		for now.After(parseDate) {
			parseDate = parseDate.AddDate(1, 0, 0)
		}

	default:
		return date, errors.New("Недопустимый символ")
	}

	out := parseDate.Format(layout)
	return out, nil
}

// GetNextDate проверяет дату по формату и возвращает следующую дату
func GetNextDate(task model.Task, layout string) (string, error) {
	now := time.Now()
	if len(task.Date) == 0 {
		task.Date = now.Format(layout)
	}
	_, err := time.Parse(layout, task.Date)
	if err != nil {

		return "", err
	}

	if task.Date < now.Format(layout) {
		if task.Repeat == "" {
			task.Date = now.Format(layout)
		} else {
			nextDate, err := NextDate(now, task.Date, task.Repeat, layout)
			if err != nil {
				return "", err
			}
			task.Date = nextDate
		}
	}
	return task.Date, nil
}
