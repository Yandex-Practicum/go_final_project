// package tasks provides tools for working on scheduler tasks
package tasks

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"time"
)

// dailyPattern takes repetition rule in "d" format, task date and now time and return repetition date of a task such as an error
func dailyPattern(now time.Time, startDate time.Time, repeat string) (string, error) {
	// Parse numeric value from repetition pattern
	days, err := strconv.Atoi(strings.TrimPrefix(repeat, "d "))
	if err != nil {
		log.Printf("Error parsing pattern value as int: %s\n", err)
		return "", err
	}
	// Define repetition range
	if days <= 0 || days > 400 {
		log.Printf("Invalid repetition range for 'd' pattern: %s\n", err)
		err = errors.New("invalid repetition range for 'd' pattern")
		return "", err
	}

	nextDate := startDate

	// Calculate repetition date if task date in the future
	nextDate = nextDate.AddDate(0, 0, days)

	// Calculate repetition date if task date in the past
	for now.After(nextDate) || nextDate == now {
		nextDate = nextDate.AddDate(0, 0, days)
	}

	return nextDate.Format("20060102"), nil
}

// yearlyPattern takes repetition rule in "y" pattern, task date and now time and return repetition date of a task such as an error
func yearlyPattern(now time.Time, startDate time.Time) (string, error) {
	// Calculate repetition date if task date in the future
	nextDate := startDate.AddDate(1, 0, 0)
	// Calculate repetition date if task date in the past
	for now.After(nextDate) || nextDate == now {
		nextDate = nextDate.AddDate(1, 0, 0)
	}
	return nextDate.Format("20060102"), nil
}

// NextDate takes repetition rule, task date as string and now time and return repetition date of a task such as an error
func NextDate(now time.Time, date string, repeat string) (string, error) {
	// Parse task date
	startDate, err := time.Parse("20060102", date)
	if err != nil {
		log.Printf("Task date is not in valid format: %s", err)
		return "", err
	}

	// Chose suitable calculation func
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
