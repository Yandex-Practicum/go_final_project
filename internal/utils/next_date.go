package utils

import (
	"errors"
	"slices"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, dateStr string, repeat string) (string, error) {
	date, err := time.Parse(ParseDateFormat, dateStr)
	if err != nil {
		return "", errors.New("неверный формат даты")
	}

	if repeat == "" {
		return "", errors.New("правило повтора не указано")
	}

	switch {
	case repeat == "y":
		for {
			date = date.AddDate(1, 0, 0)
			if date.After(now) {
				break
			}
		}
	case strings.HasPrefix(repeat, "d "):
		days, err := strconv.Atoi(strings.TrimPrefix(repeat, "d "))
		if err != nil || days < 1 || days > 400 {
			return "", errors.New("неверный 'd' формат повтора")
		}
		for {
			date = date.AddDate(0, 0, days)
			if date.After(now) {
				break
			}
		}
	case strings.HasPrefix(repeat, "w "):
		nowWeekDay := int(now.Weekday())
		if nowWeekDay == 0 {
			nowWeekDay = 7
		}
		repeatDaysStr := strings.Split(strings.TrimPrefix(repeat, "w "), ",")
		repeatDays := make([]int, 0, len(repeatDaysStr))
		for _, day := range repeatDaysStr {
			if dayNumber, parseErr := strconv.ParseInt(day, 10, 64); parseErr == nil {
				if dayNumber < 1 || dayNumber > 7 {
					return "", errors.New("неверный формат повтора")
				}
				if int(dayNumber) <= nowWeekDay {
					dayNumber += 7
				}
				repeatDays = append(repeatDays, int(dayNumber))
			}
		}
		slices.Sort(repeatDays)
		shift := repeatDays[0] - nowWeekDay
		date = now.AddDate(0, 0, shift)
	case strings.HasPrefix(repeat, "m "):
		format := strings.Split(strings.TrimPrefix(repeat, "m "), " ")
		allowDays, err := parseDays(format)
		if err != nil {
			return "", errors.New("неверный формат повтора")
		}
		allowMonths, err := parseMonths(format)
		if err != nil {
			return "", errors.New("неверный формат повтора")
		}

		for {
			if !isSliceHas(allowMonths, int(date.Month())) {
				date = date.AddDate(0, 1, 0)
				if date.Day() > 1 {
					date = date.AddDate(0, 0, -date.Day()+1)
				}
				continue
			}

			allowDaysInMonth := makeAllowDaysForMonth(date, allowDays)
			currentMonth := date.Month()
			for {
				if currentMonth != date.Month() {
					break
				}
				if isSliceHas(allowDaysInMonth, date.Day()) &&
					date.After(now) {
					return date.Format(ParseDateFormat), nil
				}
				date = date.AddDate(0, 0, 1)
			}
		}

	default:
		return "", errors.New("неверный формат повтора")
	}

	return date.Format(ParseDateFormat), nil
}

func parseDays(format []string) ([]int, error) {
	daysStr := strings.Split(format[0], ",")
	allowDays := make([]int, 0, len(daysStr))
	for _, dayS := range daysStr {
		if day, err := strconv.ParseInt(dayS, 10, 64); err == nil {
			if day < -2 || day > 31 {
				return []int{}, errors.New("неверный формат повтора")
			}
			allowDays = append(allowDays, int(day))
		}
	}

	return allowDays, nil
}

func parseMonths(format []string) ([]int, error) {
	allowMonth := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	if len(format) < 2 {
		return allowMonth, nil
	}

	monthStr := strings.Split(format[1], ",")
	allowMonth = make([]int, 0, len(monthStr))
	for _, ms := range monthStr {
		if month, err := strconv.ParseInt(ms, 10, 64); err == nil {
			if month < 1 || month > 12 {
				return []int{}, errors.New("неверный формат повтора")
			}
			allowMonth = append(allowMonth, int(month))
		}
	}
	return allowMonth, nil
}

func isSliceHas(s []int, v int) bool {
	for _, e := range s {
		if e == v {
			return true
		}
	}
	return false
}

func makeAllowDaysForMonth(date time.Time, days []int) []int {
	daysInMonth := daysIn(date.Month(), date.Year())
	result := make([]int, 0, len(days))
	for _, d := range days {
		if d > daysInMonth {
			continue
		}
		if d > 0 {
			result = append(result, d)
			continue
		}
		result = append(result, daysInMonth+d+1)
	}
	return result
}

func daysIn(m time.Month, year int) int {
	return time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}
