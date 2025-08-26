package dates

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"go_final_project/constants"
)

func dailyPattern(now time.Time, startDate time.Time, repeat string) (string, error) {

	days, err := strconv.Atoi(strings.TrimPrefix(repeat, "d "))
	if err != nil {
		log.Printf("Error parsing pattern value as int: %s\n", err)
		return "", err
	}

	if days <= 0 || days > 400 {
		log.Printf("Invalid repetition range for 'd' pattern: %s\n", err)
		err = errors.New("invalid repetition range for 'd' pattern")
		return "", err
	}

	nextDate := startDate

	nextDate = nextDate.AddDate(0, 0, days)

	for now.After(nextDate) || nextDate == now {
		nextDate = nextDate.AddDate(0, 0, days)
	}

	return nextDate.Format(constants.DateFormat), nil
}

func yearlyPattern(now time.Time, startDate time.Time) (string, error) {
	nextDate := startDate.AddDate(1, 0, 0)

	for now.After(nextDate) || nextDate == now {
		nextDate = nextDate.AddDate(1, 0, 0)
	}
	return nextDate.Format(constants.DateFormat), nil
}

func NextDate(now time.Time, date string, repeat string) (string, error) {

	startDate, err := time.Parse(constants.DateFormat, date)
	if err != nil {
		log.Printf("Task date is not in valid format: %s", err)
		return "", err
	}

	switch {
	case strings.HasPrefix(repeat, "d "):
		return dailyPattern(now, startDate, repeat)
	case repeat == "y":
		return yearlyPattern(now, startDate)
	case repeat == "":
		err = errors.New("no repetition range set")
		return "", err
	default:
		err = errors.New("repetition pattern is not supported")
		return "", err
	}
}
