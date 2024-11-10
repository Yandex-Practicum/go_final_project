package services

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, dateStr string, repeat string) (string, error) {

	date, err := time.Parse("20060102", dateStr)
	if err != nil {
		return "", errors.New("некорректный формат даты")
	}

	var nextDate time.Time

	if strings.HasPrefix(repeat, "d ") {

		daysStr := strings.TrimPrefix(repeat, "d ")
		days, err := strconv.Atoi(daysStr)
		if err != nil || days < 1 || days > 400 {
			return "", errors.New("некорректное колличество дней")
		}

		nextDate = date
		for {
			nextDate = nextDate.AddDate(0, 0, days)
			if nextDate.After(now) {
				break
			}
		}
	} else if repeat == "y" {

		nextDate = date.AddDate(1, 0, 0)
		if nextDate.Before(now) {
			nextDate = nextDate.AddDate(1, 0, 0)
		}
	} else {
		return "", errors.New("неправильный формат правила")
	}

	return nextDate.Format("20060102"), nil
}
