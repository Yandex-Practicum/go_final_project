package utils

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

func nextDaily(now, date time.Time, repeat string) (string, error) {
	args := strings.Split(repeat, " ")
	if len(args) != 2 {
		return "", fmt.Errorf("uncorrect repeat format")
	}

	v, err := strconv.Atoi(args[1])
	if err != nil {
		return "", fmt.Errorf("uncorrect repeat format")
	}

	if v < 1 || v > 400 {
		return "", fmt.Errorf("uncorrect repeat format")
	}

	date = date.AddDate(0, 0, v)
	for now.After(date) {
		date = date.AddDate(0, 0, v)
	}

	return date.Format("20060102"), nil
}

func nextWeekendly(now, date time.Time, repeat string) (string, error) {
	args := strings.Split(repeat, " ")
	if len(args) != 2 {
		return "", fmt.Errorf("uncorrect repeat format")
	}

	sd := strings.Split(args[1], ",")
	if len(sd) > 7 {
		return "", fmt.Errorf("uncorrect repeat format")
	}

	daysOfWeek := make([]int, len(sd))
	for i, el := range sd {
		v, err := strconv.Atoi(el)
		if err != nil {
			return "", fmt.Errorf("uncorrect repeat format")
		}
		if v < 1 || v > 7 {
			return "", fmt.Errorf("uncorrect repeat format")
		}

		daysOfWeek[i] = v
		sort.Slice(daysOfWeek, func(i, j int) bool { return i < j })
	}

	if now.After(date) {
		nowWD := int(now.Weekday())
		if nowWD == 0 {
			nowWD = 7
		}
		for i, v := range daysOfWeek {
			if v >= nowWD {
				date = now.AddDate(0, 0, v-nowWD)
				break
			}
			if i == len(daysOfWeek)-1 {
				date = now.AddDate(0, 0, 7-nowWD+daysOfWeek[0])
			}
		}
	}

	return date.Format("20060102"), nil
}

func nextMonthly(now, date time.Time, repeat string) (string, error) {
	return "", fmt.Errorf("uncorrect repeat format")
}

func nextYearly(now, date time.Time) (string, error) {
	date = date.AddDate(1, 0, 0)
	for now.After(date) {
		date = date.AddDate(1, 0, 0)
	}

	return date.Format("20060102"), nil
}

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("repeat is empty")
	}

	d, err := time.Parse("20060102", date)
	if err != nil {
		return "", fmt.Errorf("cannot parse date")
	}

	switch {
	case strings.HasPrefix(repeat, "d"):
		return nextDaily(now, d, repeat)
	case strings.HasPrefix(repeat, "w"):
		return "", fmt.Errorf("uncorrect repeat format")
		//return nextWeekendly(now, d, repeat)
	case strings.HasPrefix(repeat, "m"):
		return nextMonthly(now, d, repeat)
	case repeat == "y":
		return nextYearly(now, d)
	}

	return "", fmt.Errorf("unexpected type")
}
