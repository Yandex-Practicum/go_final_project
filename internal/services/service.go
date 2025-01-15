package services

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	// Преобразуем строку даты в формат time.Time
	startDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", err
	}

	// Обработка правил повторения
	var result time.Time
	var daysToAdd int
	switch {
	case repeat == "":
		// Если правило не указано, удаляем задачу
		return "", nil
	case strings.HasPrefix(repeat, "d "):
		// Преобразуем правило в число дней
		daysToAdd, err = strconv.Atoi(repeat[2:])
		if err != nil || daysToAdd > 400 || daysToAdd < 1 {
			return "", fmt.Errorf("некорректное количество дней: %s", repeat)
		}
		result = addDays(startDate, daysToAdd, now)

	case repeat == "y":
		result = addYear(startDate, now)
	// TODO добавить правила для недель и месяцев
	default:
		return "", fmt.Errorf("неподдерживаемый формат правила повторения: %s", repeat)
	}
	return result.Format("20060102"), nil

}

// addDays вычисляет следующую дату в цикле
func addDays(date time.Time, daysToAdd int, now time.Time) time.Time {
	step := daysToAdd
	for {
		nextDate := date.AddDate(0, 0, daysToAdd)
		if nextDate.After(now) {
			return nextDate
		}
		daysToAdd = daysToAdd + step
	}
}

func addYear(date time.Time, now time.Time) time.Time {
	step := 1
	for {
		nextDate := date.AddDate(step, 0, 0)
		if nextDate.After(now) {
			return nextDate
		}
		step++
	}
}
