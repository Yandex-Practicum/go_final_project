package scheduler

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("В колонке repeat — пустая строка")
	}
	startDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", errors.New("Время в переменной date не может быть преобразовано в корректную дату")
	}

	switch {
	case strings.HasPrefix(repeat, "d "):
		// d <number>
		daysStr := strings.TrimSpace(strings.TrimPrefix(repeat, "d "))
		days, err := strconv.Atoi(daysStr)
		if err != nil || days <= 0 {
			return "", errors.New("Не указан интервал в днях")
		}
		if days > 400 {
			return "", errors.New("Превышен максимально допустимый интервал")
		}
		// Если startDate равен сегодняшнему дню или находится в будущем
		if !startOfDay(startDate).Before(startOfDay(now)) {
			startDate = startDate.AddDate(0, 0, days)
		} else {
			// Если startDate в прошлом, прибавляем дни до будущего
			for startOfDay(startDate).Before(startOfDay(now)) {
				startDate = startDate.AddDate(0, 0, days)
			}
		}

	case repeat == "y":
		// y (yearly)
		startDate = startDate.AddDate(1, 0, 0)
		for !startDate.After(now) {
			startDate = startDate.AddDate(1, 0, 0)
		}

	case strings.HasPrefix(repeat, "w "):
		// w <days of week>
		daysOfWeek := strings.Split(strings.TrimSpace(strings.TrimPrefix(repeat, "w ")), ",")
		allowedDays := map[int]bool{}
		for _, dayStr := range daysOfWeek {
			day, err := strconv.Atoi(dayStr)
			if err != nil || day < 1 || day > 7 {
				return "", errors.New("Недопустимое значение")
			}
			allowedDays[day] = true
		}

		if !startDate.After(now) {
			startDate = now
		}
		startDate = startDate.AddDate(0, 0, 1)
		for {
			currentDay := int(startDate.Weekday())
			if currentDay == 0 {
				currentDay = 7
			}

			if allowedDays[currentDay] {
				break
			}

			startDate = startDate.AddDate(0, 0, 1)
		}

	case strings.HasPrefix(repeat, "m "):
		// m <days> [months]
		parts := strings.Split(strings.TrimSpace(strings.TrimPrefix(repeat, "m ")), " ")
		if len(parts) < 1 {
			return "", errors.New("Недопустимое значение")
		}

		days := strings.Split(parts[0], ",")
		allowedDays := map[int]bool{}
		for _, d := range days {
			if d == "-1" {
				allowedDays[-1] = true
				continue
			}
			if d == "-2" {
				allowedDays[-2] = true
				continue
			}
			day, err := strconv.Atoi(d)
			if err != nil || day < 1 || day > 31 {
				return "", errors.New("Недопустимый день месяца")
			}
			allowedDays[day] = true
		}

		allowedMonths := map[int]bool{}
		if len(parts) > 1 {
			months := strings.Split(parts[1], ",")
			for _, m := range months {
				month, err := strconv.Atoi(m)
				if err != nil || month < 1 || month > 12 {
					return "", errors.New("Недопустимый месяц")
				}
				allowedMonths[month] = true
			}
		} else {
			for i := 1; i <= 12; i++ {
				allowedMonths[i] = true
			}
		}

		for {
			currentDay := startDate.Day()
			currentMonth := int(startDate.Month())
			lastDayOfMonth := time.Date(startDate.Year(), startDate.Month()+1, 0, 0, 0, 0, 0, startDate.Location()).Day()

			if allowedMonths[currentMonth] {
				if allowedDays[currentDay] && startDate.After(now) {
					break
				}
				if allowedDays[-1] && currentDay == lastDayOfMonth && startDate.After(now) {
					break
				}
				if allowedDays[-2] && currentDay == lastDayOfMonth-1 && startDate.After(now) {
					break
				}
			}

			startDate = startDate.AddDate(0, 0, 1)
		}

	default:
		return "", errors.New("Указан неверный формат repeat")
	}

	return startDate.Format("20060102"), nil
}

func startOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}
