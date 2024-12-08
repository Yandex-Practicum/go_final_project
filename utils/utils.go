package utils

import (
	"errors"
	"time"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	parsedDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", err
	}

	switch {
	case repeat == "":
		return "", errors.New("no repeat rule specified")
	case repeat == "y":
		nextDate := parsedDate.AddDate(1, 0, 0)
		if nextDate.Before(now) {
			return "", errors.New("next date is before now")
		}
		return nextDate.Format("20060102"), nil
	case len(repeat) > 2 && repeat[:2] == "d ":
		days, err := time.ParseDuration(repeat[2:] + "h")
		if err != nil || days.Hours() > 400 {
			return "", errors.New("invalid repeat rule")
		}
		nextDate := parsedDate.Add(days)
		if nextDate.Before(now) {
			return "", errors.New("next date is before now")
		}
		return nextDate.Format("20060102"), nil
	default:
		return "", errors.New("unsupported repeat rule")
	}
}
