package services

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// DateFormat - формат даты, используемый в системе (YYYYMMDD).
// Используем для парсинга и форматирования дат в строковом представлении.
const DateFormat = "20060102"

// AfterNow проверяет, наступает ли дата `date` позже, чем `now`.
// Параметры:
// date - проверяемая дата.
// now - текущая дата для сравнения.
// Возвращает: true, если `date` строго больше `now` (с учётом только даты, без времени), иначе false.
func AfterNow(date, now time.Time) bool {
	// Обрезаем время до 00:00:00, чтобы сравнивать только даты (без учёта часов, минут и секунд).
	dateTruncated := date.Truncate(24 * time.Hour)
	nowTruncated := now.Truncate(24 * time.Hour)

	// Сравниваем обрезанные даты - если дата `date` после `now`, возвращаем true.
	return dateTruncated.After(nowTruncated)
}

// matchesMDay проверяет, соответствует ли дата `date` одному из указанных дней месяца.
// Параметры:
// date - проверяемая дата.
// days - список допустимых дней месяца (положительные числа 1–31, -1 - последний день месяца, -2 - предпоследний день).
// Возвращает: true, если дата соответствует одному из указанных дней, иначе false.
func matchesMDay(date time.Time, days []int) bool {
	year, month, _ := date.Date()

	// Получаем последний день месяца: создаём дату первого дня следующего месяца и вычитаем один день.
	lastDay := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()

	// Проходим по всем указанным дням из списка `days`.
	for _, day := range days {
		switch {
		// Если день в диапазоне 1-31, проверяем совпадение с днём месяца в `date`.
		case day >= 1 && day <= 31:
			if date.Day() == day {
				return true
			}
		// Если указан -1, проверяем, является ли дата последним днём месяца.
		case day == -1:
			if date.Day() == lastDay {
				return true
			}
		// Если указан -2, проверяем, является ли дата предпоследним днём месяца.
		case day == -2:
			if date.Day() == lastDay-1 {
				return true
			}
		}
	}
	return false

}

// NextDate вычисляет следующую дату по правилу повторения, начиная с `dstart`.
// Параметры:
// now - текущая дата и время (используется для сравнения).
// dstart - начальная дата в формате DateFormat (строка).
// repeat - правило повторения в виде строки (например, "d 7", "y", "w 1,2", "m 1,15 1,3,5").
// Возвращает:
// - следующую подходящую дату в формате DateFormat (строка);
// - ошибку при некорректных входных данных или невозможности вычисления даты.
func NextDate(now time.Time, dstart string, repeat string) (string, error) {

	// Парсим стартовую дату из строки в формат time.Time согласно константе DateFormat.
	date, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("failed to parse date: %w", err)
	}

	// Проверяем, что правило повторения не пустое - без правила расчёт невозможен.
	if repeat == "" {
		return "", errors.New("repeat rule is missing")
	}

	// Разбиваем правило повторения на части по пробелам для дальнейшей обработки.
	parts := strings.Split(repeat, " ")

	// Обрабатываем разные типы правил повторения (d, y, w, m).
	switch parts[0] {
	case "d":
		// Для правила "d" (дни) нужно ровно 2 части: "d" и число интервала.
		if len(parts) != 2 {
			return "", errors.New("rule 'd' requires exactly one numeric value")
		}

		// Преобразуем интервал из строки в число (количество дней).
		interval, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", fmt.Errorf("interval must be a valid integer: %w", err)
		}

		// Проверяем допустимый диапазон интервала (1-400 дней).
		if interval <= 0 || interval > 400 {
			return "", errors.New("interval must be in range [1, 400]")
		}

		// Увеличиваем дату на интервал в цикле, пока она не станет строго больше `now`.
		for {
			date = date.AddDate(0, 0, interval)
			if AfterNow(date, now) {
				break
			}
		}
	case "y":
		// Для правила "y" (год) увеличиваем дату на 1 год в цикле, пока она не превысит `now`.
		for {
			date = date.AddDate(1, 0, 0)
			if AfterNow(date, now) {
				break
			}
		}
	case "w":
		if len(parts) < 2 {
			return "", errors.New("rule 'w' requires comma-separated list of weekdays")
		}

		// Парсим дни недели из строки: разделяем по запятой и преобразуем в числа.
		dayStr := strings.Split(parts[1], ",")
		weekdays := make([]int, len(dayStr))
		for i, s := range dayStr {
			day, err := strconv.Atoi(s)
			if err != nil || day < 1 || day > 7 {
				return "", fmt.Errorf("invalid weekday value: %s", s)
			}
			// Воскресенье (7) преобразуется в 0, остальные дни - в day.
			if day == 7 {
				weekdays[i] = 0
			} else {
				weekdays[i] = day
			}
		}

		// Начинаем поиск с завтрашнего дня относительно стартовой даты.
		candidateDate := date.AddDate(0, 0, 1)

		// Увеличиваем candidateDate, пока она не станет строго больше `now`.
		for {
			candidateDate = candidateDate.AddDate(0, 0, 1)
			if AfterNow(candidateDate, now) {
				break
			}
		}

		// Ищем ближайший подходящий день недели из списка `weekdays`.
	loop:
		for {
			// Получаем номер дня недели для candidateDate (0 - воскресенье, 1 - понедельник, ..., 6 - суббота).
			weekday := int(candidateDate.Weekday())

			// Проверяем, совпадает ли текущий день недели с любым из целевых дней.
			for _, targetDay := range weekdays {
				if weekday == targetDay {
					date = candidateDate
					break loop
				}
			}

			// Если день не подошёл, переходим к следующему дню.
			candidateDate = candidateDate.AddDate(0, 0, 1)
		}

	case "m":
		if len(parts) < 2 {
			return "", errors.New("rule 'm' requires a list of days of the month")
		}

		// Парсим дни месяца из первой части правила (разделенной запятыми).
		dayPart := strings.Split(parts[1], ",")
		days := make([]int, 0, len(dayPart))

		// Преобразуем каждую строку в число и проверяем допустимость значения.
		for _, s := range dayPart {
			day, err := strconv.Atoi(s)
			if err != nil {
				return "", fmt.Errorf("day of month must be a valid integer: %s", s)
			}
			// Проверяем, что день находится в допустимом диапазоне: от -2 до 31.
			if day < -2 || day > 31 {
				return "", fmt.Errorf("day of month must be in range [-2, 31]: got %d", day)
			}
			// Добавляем корректный день в слайс days.
			days = append(days, day)
		}

		var months []int

		// Если указаны месяцы (третья часть правила), парсим их.
		if len(parts) > 2 {
			monthPart := strings.Split(parts[2], ",")

			for _, m := range monthPart {
				month, err := strconv.Atoi(m)
				if err != nil {
					return "", fmt.Errorf("month must be a valid integer: %s", m)
				}
				// Проверяем, что месяц находится в диапазоне 1–12.
				if month < 1 || month > 12 {
					return "", fmt.Errorf("month must be in range [1, 12]: got %d", month)
				}
				// Добавляем корректный месяц в срез months.
				months = append(months, month)
			}
		}

		// Начинаем поиск с завтрашнего дня относительно стартовой даты.
		candidateDate := date.AddDate(0, 0, 1)

		// Увеличиваем candidateDate, пока она не станет строго больше `now`.
		for {
			candidateDate = candidateDate.AddDate(0, 0, 1)
			if AfterNow(candidateDate, now) {
				break
			}
		}

		// Ищем ближайшую подходящую дату, соответствующую правилам дней и месяцев.
	loopTwo:
		for {
			// Получаем номер месяца для candidateDate.
			month := candidateDate.Month()
			// Получаем число дня для candidateDate.
			day := candidateDate.Day()

			// Если месяцы не указаны, проверяем только соответствие дней.
			if len(months) == 0 {
				if matchesMDay(candidateDate, days) {
					date = candidateDate
					break loopTwo
				}
			}
			// Если месяцы указаны, проверяем совпадение и месяца, и дня.
			found := false
			for _, targetMonth := range months {
				for _, targetDay := range days {
					// Если текущий месяц и день совпадают с целевыми, фиксируем дату.
					if int(month) == targetMonth && day == targetDay {
						date = candidateDate
						found = true
						break
					}
				}
				if found {
					// Выходим из обоих циклов при нахождении подходящей даты.
					break loopTwo
				}
			}
			// Если текущая дата не подошла, переходим к следующему дню.
			candidateDate = candidateDate.AddDate(0, 0, 1)
		}
	default:
		// Если правило повторения не соответствует ни одному из известных типов, возвращаем ошибку.
		return "", fmt.Errorf("unsupported repeat rule: %s", parts[0])
	}

	// Форматируем итоговую дату в требуемый строковый формат (YYYYMMDD).
	return date.Format(DateFormat), nil
}
