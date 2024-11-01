package donetaskrepeat

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
	"todo-server/internal/settings"
)

func nextWeekDay(now, startDate time.Time, daysOfWeek []int) (string, error) {
	newDate := startDate
	for {
		for _, d := range daysOfWeek {
			for newDate.Weekday() != time.Weekday(d%7) {
				newDate = newDate.AddDate(0, 0, 1)
			}
			if newDate.After(now) {
				return newDate.Format(settings.Template), nil
			}
		}
		newDate = newDate.AddDate(0, 0, 1)
	}
}

func nextMonthDate(repeat string, now, startDate time.Time) (string, error) {
	path := strings.Split(strings.TrimSpace(repeat[2:]), " ")
	if len(path) < 1 || len(path) > 2 {
		return "", fmt.Errorf("неверный формат дней и месяцев")
	}
	var days []int
	strDay := strings.Split(path[0], ",")
	for _, d := range strDay {
		day, err := strconv.Atoi(d)
		if err != nil || day == 0 || day > 31 || day < -2 {
			return "", fmt.Errorf("неверный формат дня")
		}
		days = append(days, day)
	}

	var months []int
	if len(path) == 2 {
		strMonths := strings.Split(path[1], ",")
		for _, m := range strMonths {
			month, err := strconv.Atoi(m)
			if err != nil || month > 12 || month < 1 {
				return "", fmt.Errorf("неверный формат дня")
			}
			months = append(months, month)
		}
	}

	newDate := startDate
	var dateSlice []string
	for {
		if len(months) == 0 || compareOfMonth(months, int(newDate.Month())) {
			for _, d := range days {
				dayDate := GetDay(newDate, d)
				if dayDate.After(now) {
					dateSlice = append(dateSlice, dayDate.Format(settings.Template))
				}
			}
			if len(dateSlice) != 0 {
				sort.Strings(dateSlice)
				return dateSlice[0], nil
			}
		}
		newDate = newDate.AddDate(0, 1, 0)
	}
}

func compareOfMonth(months []int, month int) bool {
	for _, m := range months {
		if m == month {
			return true
		}
	}
	return false
}

func GetDay(date time.Time, day int) time.Time {
	firstDayNextMounth := time.Date(date.Year(), date.Month()+1, 1, 0, 0, 0, 0, date.Location())
	lastDayCurrentMounth := firstDayNextMounth.AddDate(0, 0, -1).Day()
	if day > lastDayCurrentMounth {
		return time.Date(date.Year(), date.Month()+1, day, 0, 0, 0, 0, date.Location())
	}
	if day > 0 {
		return time.Date(date.Year(), date.Month(), day, 0, 0, 0, 0, date.Location())
	}
	return time.Date(date.Year(), date.Month(), lastDayCurrentMounth+day+1, 0, 0, 0, 0, date.Location())
}

func NextDate(now time.Time, date string, repeat string) (string, error) {
	startDate, err := time.Parse(settings.Template, date)
	if err != nil {
		return "", fmt.Errorf("неверный формат даты: %v", err)
	}

	if repeat == "" {
		return "", fmt.Errorf("пустая строка")
	}

	switch {
	case repeat == "y":
		newDate := startDate.AddDate(1, 0, 0)
		for newDate.Before(now) {
			newDate = newDate.AddDate(1, 0, 0)
		}
		return newDate.Format(settings.Template), nil

	case strings.HasPrefix(repeat, "d "):
		day, err := strconv.Atoi(strings.TrimSpace(repeat[2:]))
		if err != nil {
			return "", fmt.Errorf("неверный формат дня")
		}
		if day >= 1 && day <= 400 {
			newDate := startDate.AddDate(0, 0, day)
			for newDate.Before(now) {
				newDate = newDate.AddDate(0, 0, day)
			}
			return newDate.Format(settings.Template), nil
		} else {
			return "", fmt.Errorf("неверный формат дня")
		}

	case strings.HasPrefix(repeat, "w "):
		daystr := strings.Split(strings.TrimSpace(repeat[2:]), ",")
		var daysOfWeek []int
		for _, d := range daystr {
			day, err := strconv.Atoi(d)
			if err != nil || day < 1 || day > 7 {
				return "", fmt.Errorf("неверный формат дней")
			}
			daysOfWeek = append(daysOfWeek, day)
		}
		return nextWeekDay(now, startDate, daysOfWeek)

	case strings.HasPrefix(repeat, "m "):
		return nextMonthDate(repeat, now, startDate)

	default:
		return "", fmt.Errorf("неподдерживаемый формат")
	}
}
