package main

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

// NextDate вычисляет следующую дату для задачи на основе правила повторения
func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("правило повторения пустое")
	}

	// Разбираем начальную дату задачи с использованием постоянного формата DateFormat
	tDate, err := time.Parse(DateFormat, date)
	if err != nil {
		return "", fmt.Errorf("неверный формат даты: %v", err)
	}

	// Начинаем расчет с tDate или now, в зависимости от того, что больше
	startDate := now
	if tDate.After(now) {
		startDate = tDate
	}
func isLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}
	switch {
	case strings.HasPrefix(repeat, "d "): // ежедневное повторения "d"
		daysStr := strings.TrimSpace(repeat[2:])
		days, err := strconv.Atoi(daysStr)
		if err != nil || days <= 0 || days > 400 {
			return "", fmt.Errorf("неверное правило повторения 'd': %v", err)
		}
		tDate = tDate.AddDate(0, 0, days)
		for tDate.Before(now) {
			tDate = tDate.AddDate(0, 0, days)
		}
		return tDate.Format(DateFormat), nil

	case repeat == "y": // ежегодное повторения "y"
		for !tDate.After(startDate) {
			year := tDate.Year() + 1
			month := tDate.Month()
			day := tDate.Day()

			// високосный год
			if month == time.February && day == 29 && !isLeapYear(year) {
				tDate = time.Date(year, time.March, 1, 0, 0, 0, 0, tDate.Location())
			} else {
				tDate = time.Date(year, month, day, 0, 0, 0, 0, tDate.Location())
			}
		}

		return tDate.Format(DateFormat), nil

	case strings.HasPrefix(repeat, "w "): // Еженедельное повторения
		daysStr := strings.TrimSpace(repeat[2:])
		days := strings.Split(daysStr, ",")
		if len(days) == 0 {
			return "", fmt.Errorf("неверное правило повторения 'w': дни не указаны")
		}

		var daysOfWeek []int
		for _, dayStr := range days {
			day, err := strconv.Atoi(dayStr)
			if err != nil || day < 1 || day > 7 {
				return "", fmt.Errorf("неверное правило повторения 'w': неверный день '%s'", dayStr)
			}
			if day == 7 {
				day = 0
			}
			daysOfWeek = append(daysOfWeek, day)
		}
		sort.Ints(daysOfWeek)

		// Устанавливаем начальную дату startDate на максимум из tDate и now
		startDate = tDate
		if now.After(tDate) {
			startDate = now
		}

		initialDate := tDate

		// Выполняем цикл, пока день недели не совпадет с одним из daysOfWeek
		for !containsInt(daysOfWeek, int(startDate.Weekday())) || !(startDate.YearDay() > initialDate.YearDay()) {
			startDate = startDate.AddDate(0, 0, 1)
		}
		return startDate.Format(DateFormat), nil

	case strings.HasPrefix(repeat, "m "): // Обрабатываем правило ежемесячного повторения "m <дни> [<месяцы>]"
		parts := strings.Split(strings.TrimSpace(repeat[2:]), " ")

		if len(parts) == 0 {
			return "", fmt.Errorf("неверное правило повторения 'm': дни не указаны")
		}

		dayParts := strings.Split(parts[0], ",")
		var daysOfMonth []int
		for _, dayStr := range dayParts {
			day, err := strconv.Atoi(dayStr)
			if err != nil || day == 0 || day < -2 || day > 31 {
				return "", fmt.Errorf("неверное правило повторения 'm': неверный день '%s'", dayStr)
			}
			daysOfMonth = append(daysOfMonth, day)
		}

		var months []int
		if len(parts) > 1 {
			monthParts := strings.Split(parts[1], ",")
			for _, monthStr := range monthParts {
				month, err := strconv.Atoi(monthStr)
				if err != nil || month < 1 || month > 12 {
					return "", fmt.Errorf("неверное правило повторения 'm': неверный месяц '%s'", monthStr)
				}
				months = append(months, month)
			}
		}

		sort.Ints(daysOfMonth)
		sort.Ints(months)

		for {
			currentYear, currentMonth := tDate.Year(), tDate.Month()

			if len(months) > 0 && !containsInt(months, int(currentMonth)) {
				nextMonth := findNextMonth(int(currentMonth), months)
				if nextMonth <= int(currentMonth) {
					tDate = time.Date(currentYear+1, time.Month(nextMonth), 1, 0, 0, 0, 0, tDate.Location())
				} else {
					tDate = time.Date(currentYear, time.Month(nextMonth), 1, 0, 0, 0, 0, tDate.Location())
				}
				continue
			}

			// Находим ближайшую допустимую дату в текущем месяце
			var nextValidDate time.Time
			for _, day := range daysOfMonth {
				var cDate time.Time
				lastDay := lastDayOfMonth(currentYear, currentMonth)

				if day > 0 {
					if day <= lastDay {
						cDate = time.Date(currentYear, currentMonth, day, 0, 0, 0, 0, tDate.Location())
					} else {
						continue
					}
				} else {
					if -day <= lastDay {
						cDate = time.Date(currentYear, currentMonth, lastDay+day+1, 0, 0, 0, 0, tDate.Location())
					} else {
						continue
					}
				}
				// Находим ближайшую дату после startDate
				if cDate.After(startDate) && (nextValidDate.IsZero() || cDate.Before(nextValidDate)) {
					nextValidDate = cDate
				}
			}
			// Если найдена допустимая дата, возвращаем её
			if !nextValidDate.IsZero() {
				return nextValidDate.Format(DateFormat), nil
			}
			// Если не найдено допустимой даты в текущем месяце, переходим к следующему месяцу
			tDate = tDate.AddDate(0, 1, 0)
			tDate = time.Date(tDate.Year(), tDate.Month(), 1, 0, 0, 0, 0, tDate.Location())
		}

	default:
		return "", fmt.Errorf("неподдерживаемый формат правила повторения: '%s'", repeat)
	}
}

func lastDayOfMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func containsInt(slice []int, value int) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func findNextMonth(currentMonth int, months []int) int {
	for _, month := range months {
		if month >= currentMonth {
			return month
		}
	}
	return months[0]
}

func isLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}
