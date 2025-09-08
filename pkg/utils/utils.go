package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// AfterNow возвращает true, если date > now (игнорируя время)
func AfterNow(date, now time.Time) bool {
	d1 := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	d2 := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return d1.After(d2)
}

// NextDate вычисляет следующую дату задачи по правилу repeat
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	date, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", fmt.Errorf("неверный формат даты dstart: %w", err)
	}

	repeat = strings.TrimSpace(repeat)
	if repeat == "" {
		return "", errors.New("пустое правило repeat")
	}

	parts := strings.Fields(repeat)
	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", errors.New("не указан интервал дней для правила d")
		}
		interval, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", errors.New("интервал дней должен быть числом")
		}
		if interval < 1 || interval > 400 {
			return "", errors.New("интервал дней должен быть в диапазоне 1-400")
		}

		if AfterNow(date, now) {
			date = date.AddDate(0, 0, interval)
		} else {
			for !AfterNow(date, now) {
				date = date.AddDate(0, 0, interval)
			}
		}

	case "y":
		if len(parts) != 1 {
			return "", errors.New("неверный формат правила y")
		}

		if AfterNow(date, now) {
			year := date.Year() + 1
			month := date.Month()
			day := date.Day()
			if month == 2 && day == 29 && !isLeapYear(year) {
				month = 3
				day = 1
			}
			date = time.Date(year, month, day, 0, 0, 0, 0, date.Location())
		} else {
			for !AfterNow(date, now) {
				year := date.Year() + 1
				month := date.Month()
				day := date.Day()
				if month == 2 && day == 29 && !isLeapYear(year) {
					month = 3
					day = 1
				}
				date = time.Date(year, month, day, 0, 0, 0, 0, date.Location())
			}
		}

	case "w":
		if len(parts) != 2 {
			return "", errors.New("неверный формат правила w")
		}
		days, err := parseIntList(parts[1], 1, 7)
		if err != nil {
			return "", err
		}
		dayMap := make(map[int]bool)
		for _, d := range days {
			dayMap[d] = true
		}
		for {
			date = date.AddDate(0, 0, 1)
			weekday := int(date.Weekday())
			if weekday == 0 {
				weekday = 7
			}
			if dayMap[weekday] && AfterNow(date, now) {
				break
			}
		}

	case "m":
		if len(parts) < 2 {
			return "", errors.New("неверный формат правила m")
		}
		dayParts := strings.Split(parts[1], ",")
		var days []int
		for _, dp := range dayParts {
			d, err := strconv.Atoi(strings.TrimSpace(dp))
			if err != nil {
				return "", fmt.Errorf("неверный день месяца: %s", dp)
			}
			if d < -2 || d > 31 || d == 0 {
				return "", fmt.Errorf("день месяца %d вне диапазона [-2,31]", d)
			}
			days = append(days, d)
		}

		months := []int{}
		if len(parts) == 3 {
			months, err = parseIntList(parts[2], 1, 12)
			if err != nil {
				return "", err
			}
		}

		for {
			date = date.AddDate(0, 0, 1)
			day := date.Day()
			month := int(date.Month())

			validDay := false
			for _, d := range days {
				if d > 0 && d == day {
					validDay = true
				} else if d == -1 || d == -2 {
					last := time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, date.Location()).Day()
					if (d == -1 && day == last) || (d == -2 && day == last-1) {
						validDay = true
					}
				}
			}

			validMonth := len(months) == 0
			for _, m := range months {
				if month == m {
					validMonth = true
				}
			}

			if validDay && validMonth && AfterNow(date, now) {
				break
			}
		}

	default:
		return "", fmt.Errorf("неподдерживаемый формат repeat: %s", repeat)
	}

	return date.Format("20060102"), nil
}

// parseIntList парсит строку вида "1,2,3" в слайс чисел, проверяя диапазон
func parseIntList(s string, min, max int) ([]int, error) {
	var res []int
	parts := strings.Split(s, ",")
	for _, p := range parts {
		n, err := strconv.Atoi(strings.TrimSpace(p))
		if err != nil {
			return nil, fmt.Errorf("неверное число: %s", p)
		}
		if n < min || n > max {
			return nil, fmt.Errorf("число %d вне допустимого диапазона [%d,%d]", n, min, max)
		}
		res = append(res, n)
	}
	return res, nil
}

// isLeapYear возвращает true, если год високосный
func isLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}
