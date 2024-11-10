package services

import (
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {

	var nextDate time.Time

	dateStr, err := time.Parse("20060102", date)
	if err != nil || dateStr.Before(now) {
		return "", err
	}

	var days string
	symbDay := strings.HasPrefix(repeat, "d ")
	daysStr := strings.TrimPrefix(repeat, "d ")
	daysInt, err := strconv.Atoi(daysStr)

	switch {
	case symbDay:
		if err == nil && daysInt >= 1 && daysInt <= 400 {
			days = strconv.Itoa(daysInt)
			return days, nil
		}
		return "", err

	case repeat == "y":
		nextDate = dateStr.AddDate(1, 0, 0)
		return nextDate.Format("20060102"), nil
	}

	nextDate = dateStr
	for {
		nextDate = nextDate.AddDate(0, 0, daysInt)
		if nextDate.After(now) {
			break
		}
	}

	return nextDate.Format("20060102"), nil
}
