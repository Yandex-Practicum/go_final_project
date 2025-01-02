package helpers

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, date, repeat string) (string, error) {
	taskDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", err
	}

	if repeat == "" {
		return "", nil
	}

	switch {
	case strings.HasPrefix(repeat, "d "):
		days, err := strconv.Atoi(strings.TrimSpace(repeat[2:]))
		if err != nil {
			return "", err
		}

		if days < 1 || days > 400 {
			return "", errors.New("days must be between 1 and 400")
		}

		nextDate := taskDate.AddDate(0, 0, days)

		for !nextDate.After(now) {
			nextDate = nextDate.AddDate(0, 0, days)
		}

		return nextDate.Format("20060102"), nil

	case repeat == "y":
		nextDate := taskDate.AddDate(1, 0, 0)

		for !nextDate.After(now) {
			nextDate = nextDate.AddDate(1, 0, 0)
		}

		return nextDate.Format("20060102"), nil

	case strings.HasPrefix(repeat, "w "):
		days := make([]int, 0, 7)
		daysOfWeek := strings.Split(repeat[2:], ",")

		for _, day := range daysOfWeek {
			dayOfWeek, err := strconv.Atoi(day)
			if err != nil {
				return "", err
			}

			days = append(days, dayOfWeek)
		}

		for _, day := range days {
			difference := (day - int(taskDate.Weekday()) + 7) % 7 // TODO: maybe fix
			if difference == 0 {
				difference = 7
			}

			nextDate := taskDate.AddDate(0, 0, difference)

			if nextDate.After(now) {
				return nextDate.Format("20060102"), nil
			}
		}

	case strings.HasPrefix(repeat, "m "):
		monthsData := strings.Split(repeat[2:], " ")
		if len(monthsData) > 2 || len(monthsData) < 1 {
			return "", errors.New("wrong format for rule m")
		}

		days := make([]int, 0, 31)

		for _, day := range strings.Split(monthsData[0], ",") {
			dayNumber, err := strconv.Atoi(day)
			if err != nil {
				return "", err
			}

			if dayNumber < 1 || dayNumber > 31 {
				return "", errors.New("days must be between 1 and 31")
			}

			days = append(days, dayNumber)
		}

		months := make([]int, 0, 12)
		if len(monthsData) == 2 {
			for _, month := range strings.Split(monthsData[1], ",") {
				monthNumber, err := strconv.Atoi(month)
				if err != nil {
					return "", err
				}

				if monthNumber < 1 || monthNumber > 12 {
					return "", errors.New("months must be between 1 and 12")
				}

				months = append(months, monthNumber)
			}
		} else {
			months = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
		}

	}

	return "", errors.New("repeat rule not found")
}
