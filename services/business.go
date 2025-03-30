package services

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/kmakasheva/todo-list-project/domain"
)

const layout = "20060102"

func parseDate(date string) (time.Time, error) {
	parsedDate, err := time.Parse(layout, date)
	if err != nil {
		return time.Time{}, fmt.Errorf("error parsing date: %w", err)
	}
	return parsedDate, nil
}

func NextDate(now time.Time, date, repeat string) (string, error) {
	if repeat == "" || date == "" {
		return "", errors.New("invalid input data")
	}

	now = time.Date(2024, time.January, 26, 0, 0, 0, 0, time.UTC)

	dateTime, err := parseDate(date)
	if err != nil {
		return "", err
	}

	return getNextRepeatDate(dateTime, repeat)
}

func getNextRepeatDate(dateTime time.Time, repeat string) (string, error) {
	now := time.Date(2024, time.January, 26, 0, 0, 0, 0, time.UTC)

	switch repeat[0] {
	case 'y':
		return calculateYearly(dateTime, now), nil
	case 'd':
		return calculateDaily(dateTime, now, repeat)
	case 'w':
		return calculateWeekly(dateTime, now, repeat)
	case 'm':
		return calculateMonthly(dateTime, now, repeat)
	default:
		return "", errors.New("invalid repeat format")
	}
}

func calculateYearly(dateTime, now time.Time) string {
	dateTime = dateTime.AddDate(1, 0, 0)
	for dateTime.Before(now) {
		dateTime = dateTime.AddDate(1, 0, 0)
	}
	return dateTime.Format(layout)
}

func calculateDaily(dateTime, now time.Time, repeat string) (string, error) {
	parts := strings.Split(repeat, " ")
	if len(parts) != 2 {
		return "", errors.New("invalid daily repeat format, expected 'd N'")
	}
	days, err := strconv.Atoi(parts[1])
	if err != nil || days <= 0 || days > 366 {
		return "", errors.New("invalid daily interval")
	}
	return getNextOccurrence(dateTime, now, days), nil
}

func calculateWeekly(dateTime, now time.Time, repeat string) (string, error) {
	parts := strings.Split(repeat, " ")
	if len(parts) != 2 {
		return "", errors.New("invalid weekly repeat format, expected 'w D1,D2,...'")
	}
	daysOfWeek, err := parseIntList(parts[1])
	if err != nil {
		return "", errors.New("invalid days of the week format")
	}
	for _, w := range daysOfWeek {
		if w <= 0 || w > 7 {
			return "", errors.New("invalid weekly format")
		}
	}
	return getNextWeeklyOccurrence(dateTime, now, daysOfWeek), nil
}

func calculateMonthly(dateTime, now time.Time, repeat string) (string, error) {
	parts := strings.Split(repeat, " ")
	if len(parts) < 2 || len(parts) > 3 {
		return "", errors.New("invalid monthly repeat format")
	}
	daysOfMonth, err := parseIntList(parts[1])
	if err != nil {
		return "", errors.New("invalid days of month format")
	}
	for _, m := range daysOfMonth {
		if m < -2 || m > 31 {
			return "", errors.New("invalid days of month")
		}
	}
	var months []int
	if len(parts) == 3 {
		months, err = parseIntList(parts[2])
		if err != nil {
			return "", errors.New("invalid months format")
		}
	}
	return getNextMonthlyOccurrence(dateTime, now, daysOfMonth, months), nil
}

func getNextOccurrence(dateTime, now time.Time, interval int) string {
	dateTime = dateTime.AddDate(0, 0, interval)
	for dateTime.Before(now) {
		dateTime = dateTime.AddDate(0, 0, interval)
	}
	return dateTime.Format(layout)
}

func getNextWeeklyOccurrence(dateTime, now time.Time, daysOfWeek []int) string {
	var possibleDates []time.Time
	var closestPossibleDate time.Time
	for !contains(daysOfWeek, int(dateTime.Weekday())+1) || dateTime.Before(now) {
		dateTime = dateTime.AddDate(0, 0, 1)
	}
	possibleDates = append(possibleDates, dateTime)
	closestPossibleDate = possibleDates[0]
	for _, w := range possibleDates {
		if closestPossibleDate.After(w) {
			closestPossibleDate = w
		}
	}
	closestPossibleDate = closestPossibleDate.AddDate(0, 0, 1)
	return closestPossibleDate.Format(layout)
}

func getNextMonthlyOccurrence(dateTime, now time.Time, daysOfMonth, months []int) string {
	var closestPossibleDate time.Time

	for i := 0; i < 12; i++ {
		// Пропускаем месяцы, если есть ограничения
		if len(months) > 0 && !contains(months, int(dateTime.Month())) {
			dateTime = dateTime.AddDate(0, 1, 0)
			continue
		}

		possibleDates := getValidDays(dateTime, daysOfMonth)

		for _, d := range possibleDates {
			if d.After(now) && (closestPossibleDate.IsZero() || d.Before(closestPossibleDate)) {
				closestPossibleDate = d
			}
		}

		if !closestPossibleDate.IsZero() {
			return closestPossibleDate.Format(layout)
		}

		dateTime = dateTime.AddDate(0, 1, 0)
	}

	return ""
}

func getValidDays(date time.Time, daysOfMonth []int) []time.Time {
	var dates []time.Time
	lastDay := lastDayOfMonth(date)

	for _, d := range daysOfMonth {
		day := d
		if d == -1 {
			day = lastDay
		} else if d == -2 {
			day = lastDay - 1
		}

		if day >= 1 && day <= lastDay {
			dates = append(dates, time.Date(date.Year(), date.Month(), day, 0, 0, 0, 0, date.Location()))
		}
	}

	return dates
}

func lastDayOfMonth(t time.Time) int {
	t = t.AddDate(0, 1, 0)
	t = t.AddDate(0, 0, -t.Day())
	return t.Day()
}

func parseIntList(input string) ([]int, error) {
	parts := strings.Split(input, ",")
	var result []int
	for _, part := range parts {
		val, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			return nil, err
		}
		result = append(result, val)
	}
	return result, nil
}

func contains(slice []int, value int) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func ValidateTask(t domain.Task) error {
	if _, err := strconv.Atoi(t.ID); err != nil {
		return errors.New("invalid ID")
	}
	if _, err := time.Parse(layout, t.Date); err != nil {
		return errors.New("invalid date format")
	}
	if t.Title == "" {
		return errors.New("title cannot be empty")
	}
	if t.Repeat == "" {
		return errors.New("repeat field is empty")
	}
	if !strings.ContainsAny(t.Repeat[:1], "ymwd") {
		return errors.New("invalid repeat format")
	}
	return nil
}
