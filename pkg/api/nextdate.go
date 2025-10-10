package api

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

// NextDate вычисляет следующую дату выполнения задач
func NextDate(now time.Time, dateStr string, repeat string) (string, error) {
	// Проверяем параметр repeat на пустоту
	if repeat == "" {
		return "", fmt.Errorf("пустое правило повторения")
	}

	// Парсим исходную дату
	date, err := time.Parse(dateFormat, dateStr)
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
	return date.Format(dateFormat), nil
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
	return date.Format(dateFormat), nil
}

// calculateWeekly вычисляет следующую дату для еженедельного повторения
func calculateWeekly(startDate, now time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", fmt.Errorf("неверный формат для еженедельного повторения")
	}

	// Парсим дни недели (1 - понедельник, 7 - воскресенье)
	daysStr := strings.Split(parts[1], ",")
	days := make(map[int]bool)
	for _, dayStr := range daysStr {
		day, err := strconv.Atoi(dayStr)
		if err != nil || day < 1 || day > 7 {
			return "", fmt.Errorf("неверный день недели: %s", dayStr)
		}
		days[day] = true
	}

	date := startDate
	// Начинаем поиск со следующего дня
	date = date.AddDate(0, 0, 1)

	// Ищем ближайший подходящий день недели
	for i := 0; i < 365; i++ { // Ограничим поиск годом
		weekday := int(date.Weekday())
		// Преобразуем: Sunday=0, Monday=1, ... Saturday=6
		// в наш формат: Monday=1, Sunday=7
		if weekday == 0 {
			weekday = 7
		}

		if days[weekday] && afterNow(date, now) {
			return date.Format(dateFormat), nil
		}
		date = date.AddDate(0, 0, 1)
	}

	return "", fmt.Errorf("не удалось найти подходящую дату в течении года")

}

// calculateMonthly вычисляет дату для ежемесячного повторения
func calculateMonthly(startDate, now time.Time, parts []string) (string, error) {
	if len(parts) < 2 {
		return "", fmt.Errorf("неверный формат для ежмесячного повтореня")
	}

	// Парсим дни месяца
	daysStr := strings.Split(parts[1], ",")
	days := make([]int, 0)
	for _, dayStr := range daysStr {
		day, err := strconv.Atoi(dayStr)
		if err != nil || day < -2 || day == 0 || day > 31 {
			return "", fmt.Errorf("неверный день месяца: %s", dayStr)
		}
		days = append(days, day)
	}

	// Парсим месяцы (если указаны)
	months := make(map[int]bool)
	if len(parts) > 2 {
		monthsStr := strings.Split(parts[2], ",")
		for _, monthStr := range monthsStr {
			month, err := strconv.Atoi(monthStr)
			if err != nil || month < 1 || month > 12 {
				return "", fmt.Errorf("неверный месяц: %d", month)
			}
			months[month] = true
		}
	}

	date := startDate
	// Начинаем поиск со следующего дня
	date = date.AddDate(0, 0, 1)

	// Ищем ближайшую дату
	for i := 0; i < 366; i++ { // Ограничим поиск годом + 1 день
		currentMonth := int(date.Month())
		currentDay := date.Day()
		lastDay := getLastDayOfMonth(date.Year(), date.Month())

		// Проверяе5м подходит ли месяц (если указаны конкретные месяцы)
		validMonth := len(months) == 0 || months[currentMonth]

		if validMonth {
			// Проверяем все возможные дни
			for _, day := range days {
				var validDay bool
				switch {
				case day > 0:
					// Обычный день месяца
					validDay = currentDay == day
				case day == -1:
					// Последний день месяца
					validDay = currentDay == lastDay
				case day == -2:
					// Предпоследний день месяца
					validDay = currentDay == lastDay-1

				}
				if validDay && afterNow(date, now) {
					return date.Format(dateFormat), nil
				}
			}
		}

		date = date.AddDate(0, 0, 1)

	}
	return "", fmt.Errorf("не удалось найти подходящую дату в течении года")

}

// getLastDayOfMonth возвращает последний день месца
func getLastDayOfMonth(year int, month time.Month) int {
	// Создаем дату первого дня следующего месяца и вычитаем 1 день
	if month == 12 {
		return 31 // Декабрь всегда имеет 31 день
	}
	nextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)
	lastDay := nextMonth.AddDate(0, 0, -1)
	return lastDay.Day()
}

// afterNow проверяет, что дата больше now (игнорируя время)
func afterNow(date, now time.Time) bool {
	// Приводим к началу дня для сравнения только дат
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return date.After(now)
}
