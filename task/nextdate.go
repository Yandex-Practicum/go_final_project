package task

import (
	"fmt"
	"strconv"
	"time"
)

const formatOfDate = "20060102"

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("repeat is empty")
	}

	validDate, err := time.Parse(formatOfDate, date)
	if err != nil {
		return "", fmt.Errorf("incorrect date %v", err)
	}

	rule := string(repeat[0])
	rightLen := len(repeat) > 2
	var result string

	switch {
	//задача переносится на указанное число дней
	case rule == "d" && rightLen:
		result, err = everyDay(now, validDate, repeat[2:])
	// задача выполняется ежегодно
	case rule == "y":
		result, err = everyYear(now, validDate)
	default:
		return "", fmt.Errorf("incorrect repetition rule %v", err)
	}

	return result, err
}

func everyDay(now, date time.Time, days string) (string, error) {
	d, err := strconv.Atoi(days)
	if err != nil || d > 400 || d < 0 {
		return "", fmt.Errorf(`incorrect repetition rule in "d"`)
	}

	resultDate := date.AddDate(0, 0, d)
	for resultDate.Before(now) {
		resultDate = resultDate.AddDate(0, 0, d)
	}

	return resultDate.Format(formatOfDate), nil
}

func everyYear(now, date time.Time) (string, error) {
	if date.Before(now) {
		for date.Before(now) {
			date = date.AddDate(1, 0, 0)
		}
	} else {
		date = date.AddDate(1, 0, 0)
	}

	return date.Format(formatOfDate), nil
}
