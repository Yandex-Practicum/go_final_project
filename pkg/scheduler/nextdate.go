package scheduler

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// NextDate вычисляет следующую дату выполнения задачи
func NextDate(now time.Time, date string, repeat string) (string, error) {
	parsedDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", fmt.Errorf("Неверный формат даты %d", date)
	}
	//if date == now.Format("20060102") {
	//	return date, nil
	//}
	switch {
	case strings.HasPrefix(repeat, "d "):
		parts := strings.Split(repeat, " ")
		if len(parts) != 2 {
			return "", errors.New("неверный формат d")
		}

		days, err := strconv.Atoi(parts[1])
		if err != nil || days < 1 || days > 400 {
			return "", errors.New("неверное значение интервала дней")
		}

		// Рассчитываем следующую дату
		nextDate := parsedDate
		for {
			nextDate = nextDate.AddDate(0, 0, days)
			if nextDate.After(now) {
				break
			}
		}
		return nextDate.Format("20060102"), nil

	case repeat == "y":
		// Всегда добавляем год хотя бы один раз
		parsedDate = parsedDate.AddDate(1, 0, 0)
		// Продолжаем добавлять, пока не превысим now
		for !parsedDate.After(now) {
			parsedDate = parsedDate.AddDate(1, 0, 0)
		}

 case strings.HasPrefix(repeat, "w "):
     parts := strings.Split(repeat, " ")
     if len(parts) != 2 {
         return "", errors.New("Неверный формат w")
     }
    
     daysOfWeek, err := parseDaysOfWeek(parts[1])
     if err != nil {
         return "", err
     }

     // Начинаем поиск со следующего дня
     parsedDate = parsedDate.AddDate(0, 0, 1)
     for {
         if parsedDate.After(now) && containsDayOfWeek(daysOfWeek, parsedDate.Weekday()) {
             break
         }
         parsedDate = parsedDate.AddDate(0, 0, 1)
     }

	case strings.HasPrefix(repeat, "m "):
		parts := strings.Split(repeat, " ")
		if len(parts) < 2 {
			return "", errors.New("неверный формат m")
		}

		daysOfMonth, months, err := parseDaysAndMonths(parts[1:])
		if err != nil {
			return "", err
		}

		// Начинаем поиск со следующего дня
		parsedDate = parsedDate.AddDate(0, 0, 1)
		for {
			if parsedDate.After(now) &&
				containsDayOfMonth(daysOfMonth, parsedDate.Day(), parsedDate.Month()) &&
				containsMonth(months, parsedDate.Month()) {
				break
			}
			parsedDate = parsedDate.AddDate(0, 0, 1)
		}

	default:
		return "", errors.New("Неверный формат повтора")
	}

	return parsedDate.Format("20060102"), nil
}

// parseDaysOfWeek преобразует строку дней недели в массив Weekday
func parseDaysOfWeek(input string) ([]time.Weekday, error) {
	days := strings.Split(input, ",")
	result := make([]time.Weekday, 0, len(days))

	for _, day := range days {
		day = strings.TrimSpace(day) // Удаление пробелов из строки
		dayInt, err := strconv.Atoi(day)
		if err != nil || dayInt < 1 || dayInt > 7 {
			return nil, errors.New("Неверное значение дня недели")
		}
		result = append(result, time.Weekday(dayInt-1))
	}

	return result, nil
}

// parseDaysAndMonths преобразует строку дней и месяцев в массивы int
func parseDaysAndMonths(parts []string) ([]int, []time.Month, error) {
	if len(parts) == 0 {
		return nil, nil, errors.New("пропущен параметр дней")
	}

	days := strings.Split(parts[0], ",")
	daysOfMonth := make([]int, 0, len(days))

	for _, day := range days {
		dayInt, err := strconv.Atoi(day)
		if err != nil || dayInt < -2 || dayInt > 31 || dayInt == 0 {
			return nil, nil, errors.New("неправильное значение дня месяца")
		}
		daysOfMonth = append(daysOfMonth, dayInt)
	}

	months := make([]time.Month, 0)
	if len(parts) > 1 {
		monthParts := strings.Split(parts[1], ",")
		for _, month := range monthParts {
			monthInt, err := strconv.Atoi(month)
			if err != nil || monthInt < 1 || monthInt > 12 {
				return nil, nil, errors.New("неправильное значение месяца")
			}
			months = append(months, time.Month(monthInt))
		}
	}

	return daysOfMonth, months, nil
}

// containsDayOfWeek проверяет вхождение дней в слайс
func containsDayOfWeek(days []time.Weekday, day time.Weekday) bool {
	for _, d := range days {
		if d == day {
			return true
		}
	}
	return false
}

// containsDayOfMonth проверяет соблюдение правил для дня месяца
func containsDayOfMonth(days []int, targetDay int, targetMonth time.Month) bool {
	for _, day := range days {
		if day == -1 {
			lastDay := time.Date(0, targetMonth+1, 0, 0, 0, 0, 0, time.UTC).Day()
			if targetDay == lastDay {
				return true
			}
		} else if day == -2 {
			lastDay := time.Date(0, targetMonth+1, 0, 0, 0, 0, 0, time.UTC).Day()
			if targetDay == lastDay-1 {
				return true
			}
		} else if day == targetDay {
			return true
		}
	}
	return false
}

// containsMonth проверяет вхождение месяцев в слайс
func containsMonth(months []time.Month, month time.Month) bool {
	if len(months) == 0 {
		return true
	}
	for _, m := range months {
		if m == month {
			return true
		}
	}
	return false
}
