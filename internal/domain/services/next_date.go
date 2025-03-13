package services

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

// NextDate calculates the next occurrence of a task based on the given repeat rule.
func NextDate(now time.Time, date string, repeat string) (string, error) {
	// Convert the date string into a time object
	taskDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", errors.New("invalid date format")
	}

	// If no repeat rule is provided, return an error
	if repeat == "" {
		return "", errors.New("repeat rule is missing")
	}

	// Split the repeat rule into parts
	parts := strings.SplitN(repeat, " ", 2)
	if len(parts) < 1 {
		return "", errors.New("invalid repeat rule format")
	}

	ruleType := parts[0]
	ruleValue := ""
	if len(parts) == 2 {
		ruleValue = parts[1]
	}

	switch ruleType {
	case "d":
		days, err := strconv.Atoi(ruleValue)
		if err != nil || days <= 0 || days > 400 {
			return "", errors.New("invalid value for 'd' rule")
		}
		nextDate := taskDate.AddDate(0, 0, days)
		for nextDate.Before(now) {
			nextDate = nextDate.AddDate(0, 0, days)
		}
		return nextDate.Format("20060102"), nil

	case "y":
		nextDate := taskDate.AddDate(1, 0, 0)
		for nextDate.Before(now) {
			nextDate = nextDate.AddDate(1, 0, 0)
		}
		return nextDate.Format("20060102"), nil

	case "w":
		daysOfWeek, err := parseWeekdays(ruleValue)
		if err != nil {
			return "", err
		}
		nextDate := findNextWeekday(now, taskDate, daysOfWeek)
		return nextDate.Format("20060102"), nil

	case "m":
		daysOfMonth, months, err := parseMonthdays(ruleValue)
		if err != nil {
			return "", err
		}
		nextDate := findNextMonthday(now, taskDate, daysOfMonth, months)
		return nextDate.Format("20060102"), nil

	default:
		return "", errors.New("unsupported repeat rule format")
	}
}

// parseWeekdays converts a comma-separated list of weekdays into a slice of integers.
func parseWeekdays(value string) ([]int, error) {
	parts := strings.Split(value, ",")
	days := make([]int, len(parts))
	for i, part := range parts {
		day, err := strconv.Atoi(part)
		if err != nil || day < 1 || day > 7 {
			return nil, errors.New("invalid value for 'w' rule")
		}
		days[i] = day
	}
	return days, nil
}

// findNextWeekday finds the next valid weekday occurrence based on the given rule.
func findNextWeekday(now, taskDate time.Time, daysOfWeek []int) time.Time {
	for {
		taskDate = taskDate.AddDate(0, 0, 1)
		if contains(daysOfWeek, int(taskDate.Weekday())) && !taskDate.Before(now) {
			return taskDate
		}
	}
}

// parseMonthdays converts a string of month days into a slice of integers.
func parseMonthdays(value string) ([]int, []int, error) {
	parts := strings.Split(value, " ")
	if len(parts) > 2 {
		return nil, nil, errors.New("invalid format for 'm' rule")
	}

	days := make([]int, 0)
	for _, part := range strings.Split(parts[0], ",") {
		day, err := strconv.Atoi(part)
		if err != nil || (day < -2 || day > 31) || day == 0 {
			return nil, nil, errors.New("invalid value for 'm' rule")
		}
		days = append(days, day)
	}

	months := make([]int, 0)
	if len(parts) == 2 {
		for _, part := range strings.Split(parts[1], ",") {
			month, err := strconv.Atoi(part)
			if err != nil || month < 1 || month > 12 {
				return nil, nil, errors.New("invalid value for 'm' rule")
			}
			months = append(months, month)
		}
	}

	return days, months, nil
}

// findNextMonthday finds the next valid month day occurrence.
func findNextMonthday(now, taskDate time.Time, daysOfMonth, months []int) time.Time {
	for {
		if (len(months) == 0 || contains(months, int(taskDate.Month()))) &&
			contains(daysOfMonth, taskDate.Day()) && taskDate.After(now) {
			return taskDate
		}
		taskDate = taskDate.AddDate(0, 0, 1)
	}
}

// contains checks if a slice contains a specific integer.
func contains(slice []int, item int) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
