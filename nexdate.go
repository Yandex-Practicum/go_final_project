package main

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

// DateFormat — это постоянный формат для дат
const DateFormat = "20060102"

// NextDate вычисляет следующую дату для задачи на основе правила повторения
func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("поле пустое")
	}

	taskDate, err := time.Parse(DateFormat, date)
	if err != nil {
		return "", fmt.Errorf("неверный формат: %v", err)
	}

	startDate := maxTime(now, taskDate)

	switch {
	case strings.HasPrefix(repeat, "d "): // Ежедневное повторение
		days, err := strconv.Atoi(strings.TrimSpace(repeat[2:]))
		if err != nil || days <= 0 || days > 400 {
			return "", fmt.Errorf("неверный формат 'd': %v", err)
		}
		return findNextDailyDate(taskDate, now, days), nil

	case repeat == "y": // Ежегодное повторение
		return findNextYearlyDate(taskDate, startDate), nil

	case strings.HasPrefix(repeat, "w "): // Еженедельное повторение
		daysOfWeek := parseDays(strings.TrimSpace(repeat[2:]), 1, 7)
		return findNextValidDay(taskDate, now, daysOfWeek, 7), nil

	case strings.HasPrefix(repeat, "m "): // Ежемесячное повторение
		parts := strings.Split(strings.TrimSpace(repeat[2:]), " ")
		daysOfMonth := parseDays(parts[0], -2, 31)
		var months []int
		if len(parts) > 1 {
			months = parseDays(parts[1], 1, 12)
		}
		return findNextValidMonth(taskDate, now, daysOfMonth, months)

	default:
		return "", fmt.Errorf("Неверный формат")
	}
}

func maxTime(t1, t2 time.Time) time.Time {
	if t1.After(t2) {
		return t1
	}
	return t2
}

func findNextDailyDate(taskDate, now time.Time, days int) string {
	for taskDate.Before(now) {
		taskDate = taskDate.AddDate(0, 0, days)
	}
	return taskDate.Format(DateFormat)
}

func findNextYearlyDate(taskDate, startDate time.Time) string {
	for !taskDate.After(startDate) {
		taskDate = adjustYearlyDate(taskDate)
	}
	return taskDate.Format(DateFormat)
}

func adjustYearlyDate(taskDate time.Time) time.Time {
	year, month, day := taskDate.Date()
	newDate := time.Date(year+1, month, day, 0, 0, 0, 0, taskDate.Location())
	if newDate.Month() != month || newDate.Day() != day {
		newDate = time.Date(year+1, time.March, 1, 0, 0, 0, 0, taskDate.Location())
	}
	return newDate
}

func parseDays(dayStr string, min, max int) []int {
	parts := strings.Split(dayStr, ",")
	var days []int
	for _, part := range parts {
		day, err := strconv.Atoi(part)
		if err == nil && day >= min && day <= max {
			days = append(days, day)
		}
	}
	sort.Ints(days)
	return days
}

func findNextValidDay(taskDate, now time.Time, validDays []int, dayInterval int) string {
	startDate := maxTime(now, taskDate)
	for !containsInt(validDays, int(startDate.Weekday())) {
		startDate = startDate.AddDate(0, 0, 1)
	}
	return startDate.Format(DateFormat)
}

func findNextValidMonth(taskDate, now time.Time, daysOfMonth, months []int) (string, error) {
	for {
		currentYear, currentMonth := taskDate.Year(), taskDate.Month()

		if len(months) > 0 && !containsInt(months, int(currentMonth)) {
			nextMonth := findNextMonth(int(currentMonth), months)
			taskDate = adjustMonth(taskDate, currentYear, nextMonth)
			continue
		}

		nextValidDate := findNextValidDayInMonth(taskDate, daysOfMonth, currentYear, currentMonth)
		if nextValidDate.After(now) {
			return nextValidDate.Format(DateFormat), nil
		}

		taskDate = taskDate.AddDate(0, 1, 0)
	}
}

func findNextValidDayInMonth(taskDate time.Time, daysOfMonth []int, year int, month time.Month) time.Time {
	lastDay := lastDayOfMonth(year, month)
	for _, day := range daysOfMonth {
		if day > 0 && day <= lastDay || day < 0 && -day <= lastDay {
			nextDay := taskDate.AddDate(0, 0, day)
			if nextDay.After(taskDate) {
				return nextDay
			}
		}
	}
	return time.Time{}
}

func adjustMonth(taskDate time.Time, year, month int) time.Time {
	if month < int(taskDate.Month()) {
		return time.Date(year+1, time.Month(month), 1, 0, 0, 0, 0, taskDate.Location())
	}
	return time.Date(year, time.Month(month), 1, 0, 0, 0, 0, taskDate.Location())
}

func lastDayOfMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func containsInt(slice []int, value int) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func findNextMonth(currentMonth int, months []int) int {
	for _, month := range months {
		if month >= currentMonth {
			return month
		}
	}
	return months[0]
}
