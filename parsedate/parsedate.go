package parsedate

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	parseDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", fmt.Errorf("Bad date format %s", date)
	}

	switch {
	case strings.HasPrefix(repeat, "d "):
		part := strings.Split(repeat, " ")
		if len(part) != 2 {
			return "", errors.New("Bad format")
		}

		days, err := strconv.Atoi(part[1])
		if err != nil || days < 1 || days > 400 {
			return "", errors.New("Bad day interval")
		}

		nextDate := parseDate
		for {
			nextDate = nextDate.AddDate(0, 0, days)
			if nextDate.After(now) {
				break
			}
		}
		return nextDate.Format("20060102"), nil

	case repeat == "y":
		parseDate = parseDate.AddDate(1, 0, 0)
		for !parseDate.After(now) {
			parseDate = parseDate.AddDate(1, 0, 0)
		}
	case strings.HasPrefix(repeat, "w "):
		part := strings.Split(repeat, " ")
		if len(part) != 2 {
			return "", errors.New("Bad format")
		}
		sevenDays, err := parseSevenDays(part[1])
		if err != nil {
			return "", err
		}
		parseDate = parseDate.AddDate(0, 0, 1)

		for {
			if parseDate.After(now) && containsSevenDays(sevenDays, parseDate.Weekday()) {
				break
			}
			parseDate = parseDate.AddDate(0, 0, 1)
		}
	case strings.HasPrefix(repeat, "m "):
		part := strings.Split(repeat, " ")
		if len(part) < 2 {
			return "", errors.New("bad format")
		}
		monthDays, months, err := parseMonthDays(part[1:])
		if err != nil {
			return "", err
		}

		parseDate = parseDate.AddDate(0, 0, 1)
		for {
			if parseDate.After(now) && containsMonthDays(monthDays, parseDate.Day(), parseDate.Month()) && containsMonth(months, parseDate.Month()) {
				break
			}
			parseDate = parseDate.AddDate(0, 0, 1)
		}
	default:
		return "", errors.New("Bad format")
	}
	return parseDate.Format("20060102"), nil
}

func parseSevenDays(daysToParse string) ([]time.Weekday, error) {
	days := strings.Split(daysToParse, ",")
	result := make([]time.Weekday, 0, len(days))

	for _, day := range days {
		day = strings.TrimSpace(day)
		daysToInt, err := strconv.Atoi(day)
		if err != nil || daysToInt < 1 || daysToInt > 7 {
			return nil, errors.New("Bad week value")
		}
		result = append(result, time.Weekday(daysToInt-1))
	}
	return result, nil
}

func parseMonthDays(monthsToParse []string) ([]int, []time.Month, error) {
	if len(monthsToParse) == 0 {
		return nil, nil, errors.New("Miss day value")
	}

	days := strings.Split(monthsToParse[0], ",")
	monthDays := make([]int, 0, len(days))

	for _, day := range days {
		daysToInt, err := strconv.Atoi(day)
		if err != nil || daysToInt < -2 || daysToInt > 31 || daysToInt == 0 {
			return nil, nil, errors.New("Miss month value")
		}
		monthDays = append(monthDays, daysToInt)
	}
	months := make([]time.Month, 0)
	if len(monthsToParse) > 1 {
		monthSplit := strings.Split(monthsToParse[1], ",")
		for _, month := range monthSplit {
			monthsToInt, err := strconv.Atoi(month)
			if err != nil || monthsToInt < 1 || monthsToInt > 12 {
				return nil, nil, errors.New("Bad month value")
			}
			months = append(months, time.Month(monthsToInt))
		}
	}
	return monthDays, months, nil
}

func containsSevenDays(days []time.Weekday, day time.Weekday) bool {
	for _, d := range days {
		if d == day {
			return true
		}
	}
	return false
}

func containsMonthDays(days []int, trueDay int, trueMonth time.Month) bool {
	for _, day := range days {
		if day == -1 {
			lastDay := time.Date(0, trueMonth+1, 0, 0, 0, 0, 0, time.UTC).Day()
			if trueDay == lastDay {
				return true
			}

		} else if day == -2 {
			lastDay := time.Date(0, trueMonth+1, 0, 0, 0, 0, 0, time.UTC).Day()
			if trueDay == lastDay-1 {
				return true
			}
		} else if day == trueDay {
			return true
		}
	}
	return false
}

func containsMonth(months []time.Month, month time.Month) bool {
	if len(months) == 0 {
		return true
	}

	for _, i := range months {
		if i == month {
			return true
		}
	}
	return false
}
