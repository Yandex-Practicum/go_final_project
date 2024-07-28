package task

import (
	"fmt"
	"strconv"
	"time"
)

const timeFormat = "20060102"

func NextDate(now time.Time, date string, repeat string) (string, error) {
	var res string

	if repeat == "" {
		return "", fmt.Errorf("поле repeat пустое")
	}

	startDate, err := time.Parse(timeFormat, date)
	if err != nil {
		return "", fmt.Errorf("Некорретная дата %v", err)
	}

	rule := string(repeat[0])
	ruleLen := len(repeat) > 2

	switch {
	case rule == "d" && ruleLen:
		res, err = everyDay(now, startDate, repeat[2:])
	case rule == "y":
		res, err = everyYear(now, startDate)
	default:
		return "", fmt.Errorf("неподдерживаемый формат %v", err)
	}

	return res, err
}

func everyDay(now time.Time, date time.Time, days string) (string, error) {
	day, err := strconv.Atoi(days)
	if err != nil || day > 400 || day < 0 {
		return "", fmt.Errorf("некорректное правило")
	}

	res := date.AddDate(0, 0, day)

	for res.Before(now) {
		res = res.AddDate(0, 0, day)
	}

	return res.Format(timeFormat), nil
}

func everyYear(now time.Time, date time.Time) (string, error) {
	if date.Before(now) {
		for date.Before(now) {
			date = date.AddDate(1, 0, 0)
		}
	} else {
		date = date.AddDate(1, 0, 0)
	}

	return date.Format(timeFormat), nil
}
