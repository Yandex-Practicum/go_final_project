package main

import (
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	startDate, err := time.Parse("20060102", date)
	if err != nil || startDate.Year() < 1600 || startDate.Year() > 9999 {
		return "", nil
	}

	parts := strings.Split(repeat, " ")
	switch parts[0] {
	case "y":
		if len(parts) > 1 {
			return "", nil
		}
		nextDate := startDate.AddDate(1, 0, 0)
		for !nextDate.After(now) {
			nextDate = nextDate.AddDate(1, 0, 0)
		}
		if startDate.Month() == time.February && startDate.Day() == 29 && !isLeapYear(nextDate.Year()) {
			nextDate = time.Date(nextDate.Year(), time.March, 1, 0, 0, 0, 0, time.UTC)
		}
		return nextDate.Format("20060102"), nil

	case "d":
		if len(parts) < 2 {
			return "", nil
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days <= 0 || days > 365 {
			return "", nil
		}
		nextDate := startDate.AddDate(0, 0, days)
		for !nextDate.After(now) {
			nextDate = nextDate.AddDate(0, 0, days)
		}
		return nextDate.Format("20060102"), nil
	}

	return "", nil
}

func isLeapYear(year int) bool {
	return (year%4 == 0 && year%100 != 0) || (year%400 == 0)
}
