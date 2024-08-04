package nextdate

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Функция для преобразования строки в формат времени
func parseTime(s string) (time.Time, error) {
	t, err := time.Parse("20060102", s)
	return t, err
}
func NextDate(now time.Time, date string, repeat string) (string, error) {
	taskDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", fmt.Errorf("invalid date format: %w", err)
	}

	if repeat == "" {
		return "", fmt.Errorf("no repeat rule specified")
	}

	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", fmt.Errorf("invalid repeat rule format")
	}

	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid daily repeat rule format")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days < 1 || days > 400 {
			return "", fmt.Errorf("invalid daily repeat interval")
		}
		for {
			taskDate = taskDate.AddDate(0, 0, days)
			if taskDate.After(now) {
				return taskDate.Format("20060102"), nil
			}
		}

	case "y":
		if len(parts) != 1 {
			return "", fmt.Errorf("invalid yearly repeat rule format")
		}
		for {
			taskDate = taskDate.AddDate(1, 0, 0)
			if taskDate.After(now) {
				return taskDate.Format("20060102"), nil
			}
		}

	case "w":
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid weekly repeat rule format")
		}
		daysOfWeek := parseDaysOfWeek(parts[1])
		if len(daysOfWeek) == 0 {
			return "", fmt.Errorf("invalid weekly repeat days")
		}
		for {
			taskDate = taskDate.AddDate(0, 0, 1)
			if taskDate.After(now) && contains(daysOfWeek, int(taskDate.Weekday())) {
				return taskDate.Format("20060102"), nil
			}
		}

	case "m":
		if len(parts) < 2 {
			return "", fmt.Errorf("invalid monthly repeat rule format")
		}
		daysOfMonth, months := parseDaysAndMonths(parts[1:])
		if len(daysOfMonth) == 0 {
			return "", fmt.Errorf("invalid monthly repeat days")
		}
		for {
			taskDate = taskDate.AddDate(0, 0, 1)
			if taskDate.After(now) && contains(daysOfMonth, taskDate.Day()) && (len(months) == 0 || contains(months, int(taskDate.Month()))) {
				return taskDate.Format("20060102"), nil
			}
		}

	default:
		return "", fmt.Errorf("unsupported repeat rule")
	}
}

func parseDaysOfWeek(s string) []int {
	parts := strings.Split(s, ",")
	days := make([]int, 0, len(parts))
	for _, part := range parts {
		day, err := strconv.Atoi(part)
		if err != nil || day < 1 || day > 7 {
			return nil
		}
		days = append(days, (day-1)%7)
	}
	return days
}

func parseDaysAndMonths(parts []string) ([]int, []int) {
	days := make([]int, 0)
	months := make([]int, 0)
	if len(parts) > 0 {
		dayParts := strings.Split(parts[0], ",")
		for _, part := range dayParts {
			day, err := strconv.Atoi(part)
			if err != nil || day == 0 || day < -31 || day > 31 {
				return nil, nil
			}
			days = append(days, day)
		}
	}
	if len(parts) > 1 {
		monthParts := strings.Split(parts[1], ",")
		for _, part := range monthParts {
			month, err := strconv.Atoi(part)
			if err != nil || month < 1 || month > 12 {
				return nil, nil
			}
			months = append(months, month)
		}
	}
	return days, months
}

func contains(slice []int, value int) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

/*
// Функция для вычисления следующей даты в соответствии с правилами повторения
func NextDate(now time.Time, date string, repeat string) (string, error) {

	// Преобразование входных параметров во время
	nowTime := now.Truncate(24 * time.Hour)
	dateTime, err := parseTime(date)
	if err != nil {
		return "", fmt.Errorf("incorrect date: %w", err)
	}

	symbols := strings.Split(repeat, " ")
	if len(symbols) == 0 {
		return "", fmt.Errorf("invalid repeat format")
	}
	firstsymbol := symbols[0]

	if !strings.ContainsAny(firstsymbol, "dy") {
		// Исправляем форматирование JSON-сообщения об ошибке
		return "", fmt.Errorf(`{"error":"incorrect symbol %s"}`, firstsymbol)
	}

	switch firstsymbol {
	case "y": // Ежегодно
		if len(symbols) != 1 {
			return "", errors.New("invalid repeat format for 'y'")
		}
		return nextYearlyDate(dateTime, nowTime)

	case "d": // Через указанное число дней
		if len(symbols) != 2 {
			// Проверяем на nil err перед использованием
			if err != nil {
				return "", fmt.Errorf("incorrect repeat rule interval in days: %w", err)
			} else {
				return "", errors.New("incorrect repeat rule interval in days")
			}
		}
		secondsymbol := symbols[1]
		return nextDayRepeat(dateTime, nowTime, secondsymbol)

	default: // Неподдерживаемые форматы
		return "", errors.New("unsupported repeat format")
	}

}

// nextYearlyDate переносит дату на один год вперед
func nextYearlyDate(dateTime time.Time, nowTime time.Time) (string, error) {
	next := dateTime

	for {
		if next.After(nowTime) && next.After(dateTime) {
			break
		}
		next = next.AddDate(1, 0, 0)
	}

	return next.Format("20060102"), nil
}

// nextDayRepeat вычисляет следующую дату на основании повторения через указанное количество дней
func nextDayRepeat(dateTime time.Time, nowTime time.Time, daysInterval string) (string, error) {
	daysIntervalInt, err := strconv.Atoi(daysInterval)
	if err != nil {
		return "", fmt.Errorf("incorrect repeat rule interval: %w", err)
	}

	if daysIntervalInt > 400 {
		return "", fmt.Errorf("400 is maximum interval of days")
	}

	nextDate := dateTime.AddDate(0, 0, daysIntervalInt)

	for {
		if nextDate.After(nowTime) {
			break
		}
		nextDate = nextDate.AddDate(0, 0, daysIntervalInt)
	}
	return nextDate.Format("20060102"), nil
}
*/
