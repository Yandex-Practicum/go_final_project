package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	// Проверка на пустую строку для repeat
	if repeat == "" {
		return "", fmt.Errorf("правило повторения не указано")
	}

	// Парсинг исходной даты
	taskDate, err := time.Parse(dateLayout, date)
	if err != nil {
		return "", fmt.Errorf("некорректная дата: %v", err)
	}

	switch {
	case repeat == "y":
		taskDate = taskDate.AddDate(1, 0, 0)
		for taskDate.Before(now) || taskDate.Equal(now) {
			taskDate = taskDate.AddDate(1, 0, 0)
		}
		return taskDate.Format(dateLayout), nil

	case strings.HasPrefix(repeat, "d "):
		days, err := strconv.Atoi(strings.TrimPrefix(repeat, "d "))
		if err != nil || days <= 0 || days > 400 {
			return "", fmt.Errorf("некорректное правило d: %v", repeat)
		}
		if days == 1 && now.After(taskDate) {
			return now.Format(dateLayout), nil
		}
		taskDate = taskDate.AddDate(0, 0, days)
		for taskDate.Before(now) || taskDate.Equal(now) {
			taskDate = taskDate.AddDate(0, 0, days)
		}
		return taskDate.Format(dateLayout), nil

	case strings.HasPrefix(repeat, "w "):
		daysOfWeek := strings.Split(strings.TrimPrefix(repeat, "w "), ",")
		return nextWeekday(taskDate, daysOfWeek)

	case strings.HasPrefix(repeat, "m "):
		return nextMonthDay(taskDate, strings.TrimPrefix(repeat, "m "))

	default:
		return "", fmt.Errorf("неизвестный формат правила повторения: %v", repeat)
	}
}

// nextWeekday вычисляет следующую дату для заданных дней недели
func nextWeekday(startDate time.Time, daysOfWeek []string) (string, error) {
	nextDate := startDate

	for _, day := range daysOfWeek {
		dayNum, err := strconv.Atoi(day)
		if err != nil || dayNum < 1 || dayNum > 7 {
			return "", fmt.Errorf("недопустимое значение дня недели: %v", day)
		}
		targetDay := time.Weekday(dayNum % 7)

		// Находим ближайшую дату
		for nextDate.Weekday() != targetDay {
			nextDate = nextDate.AddDate(0, 0, 1)
		}
		if nextDate.After(startDate) {
			return nextDate.Format("20060102"), nil
		}
	}
	return nextDate.Format("20060102"), nil
}

// nextMonthDay вычисляет следующую дату для заданных дней месяца
func nextMonthDay(startDate time.Time, monthRule string) (string, error) {
	dayRules := strings.Split(monthRule, " ")
	var days []int

	for _, rule := range dayRules {
		if strings.Contains(rule, ",") {
			for _, d := range strings.Split(rule, ",") {
				day, _ := strconv.Atoi(d)
				days = append(days, day)
			}
		} else {
			day, _ := strconv.Atoi(rule)
			days = append(days, day)
		}
	}

	nextDate := startDate
	for _, day := range days {
		if day == -1 { // Последний день месяца
			nextDate = time.Date(nextDate.Year(), nextDate.Month()+1, 1, 0, 0, 0, 0, nextDate.Location()).AddDate(0, 0, -1)
		} else if day == -2 { // Предпоследний день месяца
			nextDate = time.Date(nextDate.Year(), nextDate.Month()+1, 1, 0, 0, 0, 0, nextDate.Location()).AddDate(0, 0, -2)
		} else {
			nextDate = time.Date(nextDate.Year(), nextDate.Month(), day, 0, 0, 0, 0, nextDate.Location())
		}

		// Проверяем, что следующая дата больше текущей
		if nextDate.After(startDate) {
			return nextDate.Format("20060102"), nil
		}
	}
	return "", fmt.Errorf("нет подходящей даты для заданных правил")
}
