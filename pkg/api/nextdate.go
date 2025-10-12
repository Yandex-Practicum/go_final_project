package api

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

// NextDate вычисляет следующую дату выполнения задач
func NextDate(now time.Time, dateStr string, repeat string) (string, error) {
	// Проверяем параметр repeat на пустоту
	if repeat == "" {
		return "", fmt.Errorf("пустое правило повторения")
	}

	// Парсим исходную дату
	date, err := time.Parse(DateFormat, dateStr)
	if err != nil {
		return "", fmt.Errorf("неверный формат даты: %v", err)
	}

	// Разбиваем правило на части
	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", fmt.Errorf("неверный формат строки повторения")
	}

	switch parts[0] {
	case "d":
		// Формат: d (число)
		return calculateDaily(date, now, parts)
	case "y":
		// Правило: y (ежегодно)
		return calculateYearly(date, now, parts)
	case "w":
		// Правило: w (дни недели)
		return calculateWeekly(date, now, parts)
	case "m":
		// Правило: m (дни месяца/месяцы)
		return calculateMonthly(date, now, parts)
	default:
		// Для остальных правил возвращаем ошибку
		return "", fmt.Errorf("не поддерживаемый формат повторения")
	}
}

// calculateDaily вычисляет следующую дату для ежедневного интервала
func calculateDaily(startDate, now time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", fmt.Errorf("неверный формат интервала дней")
	}

	interval, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", fmt.Errorf("неверный интервал дней: %v", err)
	}

	if interval <= 0 || interval > 400 {
		return "", fmt.Errorf("интервал дней должен быть от 1 до 400")
	}

	date := startDate
	for {
		date = date.AddDate(0, 0, interval)
		if afterNow(date, now) {
			break
		}
	}
	return date.Format(DateFormat), nil
}

// calculateYearly вычисляет следующую дату для ежегодного повторения
func calculateYearly(startDate, now time.Time, parts []string) (string, error) {
	if len(parts) != 1 {
		return "", fmt.Errorf("неверный формат ежегодного повторения")
	}

	date := startDate
	for {
		date = date.AddDate(1, 0, 0)
		if afterNow(date, now) {
			break
		}
	}
	return date.Format(DateFormat), nil
}

// calculateWeekly вычисляет следующую дату для еженедельного повторения
func calculateWeekly(startDate, now time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", fmt.Errorf("неверный формат для еженедельного повторения")
	}

	// Парсим дни недели (1 - понедельник, 7 - воскресенье)
	daysStr := strings.Split(parts[1], ",")
	weekdays := make([]int, 0, len(daysStr))
	for _, dayStr := range daysStr {
		day, err := strconv.Atoi(dayStr)
		if err != nil || day < 1 || day > 7 {
			return "", fmt.Errorf("неверный день недели: %s", dayStr)
		}
		weekdays = append(weekdays, day)
	}

	// Начинаем поиск со следующего дня
	current := startDate.AddDate(0, 0, 1)

	// Ищем ближайший подходящий день в течение 2 лет
	for i := 0; i < 365*5; i++ {
		currentWeekday := int(current.Weekday())
		// Преобразуем: Sunday=0, Monday=1, ... Saturday=6
		// в наш формат: Monday=1, Sunday=7
		if currentWeekday == 0 {
			currentWeekday = 7
		}

		// Проверяем подходит ли текущий день недели
		for _, targetWeekday := range weekdays {
			if currentWeekday == targetWeekday {
				if afterNow(current, now) {
					return current.Format(DateFormat), nil
				}
			}
		}

		current = current.AddDate(0, 0, 1)
	}

	return "", fmt.Errorf("не удалось найти подходящую дату в течении 2 лет")
}

// calculateMonthly вычисляет дату для ежемесячного повторения
func calculateMonthly(startDate, now time.Time, parts []string) (string, error) {
	if len(parts) < 2 {
		return "", fmt.Errorf("неверный формат для ежемесячного повторения")
	}

	// Парсим дни месяца
	daysStr := strings.Split(parts[1], ",")
	monthDays := make([]int, 0, len(daysStr))
	for _, dayStr := range daysStr {
		day, err := strconv.Atoi(dayStr)
		if err != nil || day < -2 || day == 0 || day > 31 {
			return "", fmt.Errorf("неверный день месяца: %s", dayStr)
		}
		monthDays = append(monthDays, day)
	}

	// Парсим месяцы (если указаны)
	targetMonths := make(map[int]bool)
	if len(parts) > 2 {
		monthsStr := strings.Split(parts[2], ",")
		for _, monthStr := range monthsStr {
			month, err := strconv.Atoi(monthStr)
			if err != nil || month < 1 || month > 12 {
				return "", fmt.Errorf("неверный месяц: %s", monthStr)
			}
			targetMonths[month] = true
		}
	}

	// Начинаем поиск со следующего дня после startDate
	current := startDate.AddDate(0, 0, 1)

	// Ищем ближайшую дату в течение 5 лет
	for i := 0; i < 365*5; i++ {
		currentMonth := int(current.Month())
		currentYear := current.Year()
		currentDay := current.Day()

		// Проверяем подходит ли месяц (если указаны конкретные месяцы)
		if len(targetMonths) > 0 && !targetMonths[currentMonth] {
			current = current.AddDate(0, 0, 1)
			continue
		}

		// Сортируем дни для поиска от меньшего к большему
		sortedDays := make([]int, len(monthDays))
		copy(sortedDays, monthDays)
		sort.Ints(sortedDays)

		// Проверяем все возможные дни
		for _, targetDay := range sortedDays {
			var actualDay int
			lastDay := getLastDayOfMonth(currentYear, time.Month(currentMonth))

			if targetDay > 0 {
				// Обычный день месяца
				if targetDay > lastDay {
					continue // Пропускаем если дня нет в этом месяце
				}
				actualDay = targetDay
			} else if targetDay == -1 {
				// Последний день месяца
				actualDay = lastDay
			} else if targetDay == -2 {
				// Предпоследний день месяца
				if lastDay < 2 {
					continue
				}
				actualDay = lastDay - 1
			}

			// Если текущий день совпадает с целевым и дата после now
			if currentDay == actualDay && afterNow(current, now) {
				return current.Format(DateFormat), nil
			}
		}

		current = current.AddDate(0, 0, 1)
	}

	return "", fmt.Errorf("не удалось найти подходящую дату в течении 5 лет")
}

// getLastDayOfMonth возвращает последний день месяца
func getLastDayOfMonth(year int, month time.Month) int {
	if month == time.December {
		return 31
	}
	// Первый день следующего месяца минус 1 день = последний день текущего месяца
	firstOfNextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)
	lastDay := firstOfNextMonth.AddDate(0, 0, -1)
	return lastDay.Day()
}

// afterNow проверяет, что дата больше now (игнорируя время)
func afterNow(date, now time.Time) bool {
	// Приводим к началу дня для сравнения только дат
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return date.After(now)
}
