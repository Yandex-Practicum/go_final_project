package api

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var checkList = map[string]struct{}{
	"d": {},
	"w": {},
	"m": {},
	"y": {},
}

func afterNow(date, now time.Time) bool {
	newDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	newNow := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	return newDate.After(newNow)
}

func nextDate(now time.Time, dstart string, repeat string) (string, error) {
	date, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", fmt.Errorf("Could not to parse dstart %w\n", err)
	}

	meta := strings.Split(repeat, " ")
	if len(meta) == 0 {
		return "", errors.New("Empty repeat rule")
	}

	if _, ok := checkList[meta[0]]; !ok {
		return "", errors.New("Invalid symbol")
	}

	switch meta[0] {
	case "d":
		if len(meta) < 2 {
			return "", errors.New("Day interval required for 'd' rule")
		}

		day, err := checkError(meta[1], 1, 400)
		if err != nil {
			return "", err
		}

		for {
			date = date.AddDate(0, 0, day)
			if afterNow(date, now) {
				break
			}
		}
	case "w":
		if len(meta) < 2 {
			return "", errors.New("Day interval required for 'w' rule")
		}

		week := [7]bool{}
		for _, v := range strings.Split(meta[1], ",") {
			day, err := checkError(v, 1, 7)
			if err != nil {
				return "", err
			}
			week[day%7] = true
		}
		for {
			date = date.AddDate(0, 0, 1)
			day := date.Weekday()
			if week[day] && afterNow(date, now) {
				break
			}
		}
	case "m":
		if len(meta) < 2 {
			return "", errors.New("Day interval required for 'm' rule")
		}

		days := [32]bool{}
		months := [13]bool{}
		switch len(meta) {
		case 2:
			for i := range months {
				months[i] = true
			}
		case 3:
			for _, v := range strings.Split(meta[2], ",") {
				mon, err := checkError(v, 1, 12)
				if err != nil {
					return "", err
				}
				months[mon] = true
			}
		}

		year, month, _ := date.Date()
		firstOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
		lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
		daysInMonth := lastOfMonth.Day()

		for _, v := range strings.Split(meta[1], ",") {
			day, err := checkError(v, -2, 31)
			if err != nil {
				return "", err
			}
			if day == 0 {
				return "", errors.New("Invalid month value transferred")
			}

			if day < 0 {
				days[daysInMonth+1+day] = true
			} else {
				days[day] = true
			}
		}

		for {
			date = date.AddDate(0, 0, 1)
			day := date.Day()
			month := date.Month()
			if days[day] && months[month] && afterNow(date, now) {
				break
			}
		}

	case "y":
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
		}
	}

	outDate := date.Format("20060102")
	return outDate, nil
}

func checkError(v string, down, up int) (int, error) {
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("Could not parse to int from string %w\n", err)
	}
	if n < down || n > up {
		return 0, errors.New("Invalid value transferred")
	}

	return n, nil
}
