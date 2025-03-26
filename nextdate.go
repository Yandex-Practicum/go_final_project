package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, dateStr string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("пустое правило повторения")
	}

	date, err := time.Parse("20060102", dateStr)
	if err != nil {
		return "", fmt.Errorf("неверный формат даты: %v", err)
	}

	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", errors.New("неверный формат правила")
	}

	switch parts[0] {
	case "d":
		return handleDailyRule(now, date, parts)
	case "y":
		return handleYearlyRule(now, date)
	case "w":
		return handleWeeklyRule(now, date, parts)
	case "m":
		return handleMonthlyRule(now, date, parts)
	default:
		return "", errors.New("неподдерживаемый формат правила")
	}
}

func handleDailyRule(now, date time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", errors.New("неверный формат правила 'd'")
	}

	days, err := strconv.Atoi(parts[1])
	if err != nil || days <= 0 || days > 400 {
		return "", errors.New("неверное количество дней")
	}

	next := date
	for {
		next = next.AddDate(0, 0, days)
		if next.After(now) {
			break
		}
	}
	return next.Format("20060102"), nil
}

func handleYearlyRule(now, date time.Time) (string, error) {
	next := date
	for {
		next = next.AddDate(1, 0, 0)
		// Особый случай для 29 февраля
		if next.Month() == time.February && next.Day() == 29 {
			if !next.After(now) {
				// Если дата всё ещё не в будущем, переходим к следующему году
				continue
			}
			// Проверяем существует ли 29 февраля в этом году
			nextYear := next.Year()
			feb29 := time.Date(nextYear, time.February, 29, 0, 0, 0, 0, time.UTC)
			if feb29.Day() != 29 {
				// Если 29 февраля не существует, переносим на 1 марта
				next = time.Date(nextYear, time.March, 1, 0, 0, 0, 0, time.UTC)
			}
		}
		if next.After(now) {
			break
		}
	}
	return next.Format("20060102"), nil
}

// ... остальные функции (handleWeeklyRule, handleMonthlyRule) остаются без изменений ...

// Обработка правила "w <дни недели>" (со звёздочкой)
func handleWeeklyRule(now, date time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", errors.New("неверный формат правила 'w'")
	}

	days, err := parseWeekDays(parts[1])
	if err != nil {
		return "", err
	}

	next := date
	for {
		next = next.AddDate(0, 0, 1)
		if next.Before(now) {
			continue
		}

		weekday := int(next.Weekday())
		if weekday == 0 {
			weekday = 7 // Воскресенье -> 7
		}

		if contains(days, weekday) {
			return next.Format("20060102"), nil
		}
	}
}

// Обработка правила "m <дни месяца> [месяцы]" (со звёздочкой)
func handleMonthlyRule(now, date time.Time, parts []string) (string, error) {
	if len(parts) < 2 {
		return "", errors.New("неверный формат правила 'm'")
	}

	days, months, err := parseMonthDays(parts[1:]...)
	if err != nil {
		return "", err
	}

	next := date
	for {
		next = next.AddDate(0, 0, 1)
		if next.Before(now) {
			continue
		}

		day := next.Day()
		month := int(next.Month())

		// Проверка специальных дней (-1, -2)
		lastDay := time.Date(next.Year(), next.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
		specialDay := 0
		if day == lastDay {
			specialDay = -1
		} else if day == lastDay-1 {
			specialDay = -2
		}

		dayMatch := contains(days, day) || contains(days, specialDay)
		monthMatch := len(months) == 0 || contains(months, month)

		if dayMatch && monthMatch {
			return next.Format("20060102"), nil
		}
	}
}

// Вспомогательные функции
func parseWeekDays(s string) ([]int, error) {
	var days []int
	for _, part := range strings.Split(s, ",") {
		d, err := strconv.Atoi(part)
		if err != nil || d < 1 || d > 7 {
			return nil, errors.New("неверный день недели")
		}
		days = append(days, d)
	}
	return days, nil
}

func parseMonthDays(parts ...string) (days []int, months []int, err error) {
    if len(parts) == 0 {
        return nil, nil, errors.New("не указаны дни месяца")
    }

    // Парсим дни месяца (первая часть правила)
    dayParts := strings.Split(parts[0], ",")
    for _, part := range dayParts {
        d, err := strconv.Atoi(part)
        if err != nil {
            return nil, nil, fmt.Errorf("неверный день месяца: %v", part)
        }

        // Проверяем допустимые значения дней
        if d == -1 || d == -2 {
            // Специальные значения (последний и предпоследний день месяца)
            days = append(days, d)
            continue
        }

        if d < 1 || d > 31 {
            return nil, nil, fmt.Errorf("день месяца должен быть от 1 до 31, или -1/-2, получили: %d", d)
        }
        days = append(days, d)
    }

    // Если есть вторая часть - парсим месяцы
    if len(parts) > 1 {
        monthParts := strings.Split(parts[1], ",")
        for _, part := range monthParts {
            m, err := strconv.Atoi(part)
            if err != nil || m < 1 || m > 12 {
                return nil, nil, fmt.Errorf("месяц должен быть от 1 до 12, получили: %v", part)
            }
            months = append(months, m)
        }
    }

    return days, months, nil
}

func contains(slice []int, val int) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}