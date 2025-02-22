package tasks

import (
	"errors"
	"strings"
	"time"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	taskDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", err
	}

	if repeat == "" {
		return "", errors.New("invalid repeat rule")
	}
	rules := strings.Split(repeat, " ")
	switch rules[0] {
	case "d":
		return DailyRepeat(now, taskDate, rules)
	case "y":
		return YearRepeat(now, taskDate, rules)
	case "w":
		return WeekRepeat(now, taskDate, rules)
	case "m":
		return MonthRepeat(now, taskDate, rules)
	default:
		return "", errors.New("invalid repeat rule")
	}
}
