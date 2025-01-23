package rules

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

// parseWeekDays извлекает дни недели из строки.
func parseWeekDays(days string) ([]time.Weekday, error) {
	dayString := strings.Split(days, ",")
	var weekdays []time.Weekday
	for _, dayStr := range dayString {
		dayInt, err := strconv.Atoi(dayStr)
		if err != nil || dayInt < 1 || dayInt > 7 {
			return nil, errors.New("некорректный день недели")
		}
		weekday := time.Weekday((dayInt % 7)) // Преобразуем в тип time.Weekday
		weekdays = append(weekdays, weekday)
	}
	return weekdays, nil
}

// handleWeekRepeat обрабатывает повторение по неделям.
func handleWeekRepeat(now time.Time, taskDate time.Time, rules []string) (string, error) {
	if len(rules) != 2 {
		return "", errors.New("неверный формат правила повторения для недели")
	}
	weekdays, err := parseWeekDays(rules[1]) // Извлекаем подстроку после "w "
	if err != nil {
		return "", err
	}

	for {
		// Проверяем каждый день недели
		for _, weekday := range weekdays {
			if taskDate.Weekday() == weekday {
				if taskDate.After(now) {
					return taskDate.Format("20060102"), nil // Возвращаем форматированную дату
				}
			}
		}
		taskDate = taskDate.AddDate(0, 0, 1) // Увеличиваем дату на один день
	}
}
