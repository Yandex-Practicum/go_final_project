package logic

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// NextDate рассчитывает следующую дату для выполнения задачи с учетом базовых правил повторения.
func NextDate(now time.Time, date string, repeat string) (string, error) {
	// Парсим входную дату в формат времени
	startDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", fmt.Errorf("не удалось разобрать дату: %v", err)
	}

	// Проверяем, что правило повторения указано
	if repeat == "" {
		return "", errors.New("правило повторения не указано")
	}

	// Обработка правила "y" (ежегодно)
	if repeat == "y" {
		nextDate := startDate.AddDate(1, 0, 0)
		for !nextDate.After(now) {
			nextDate = nextDate.AddDate(1, 0, 0)
		}
		return nextDate.Format("20060102"), nil
	}

	// Обработка правила "d <число>"
	if strings.HasPrefix(repeat, "d ") {
		parts := strings.Fields(repeat)
		if len(parts) != 2 {
			return "", errors.New("неверный формат правила повторения")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days <= 0 || days > 400 {
			return "", errors.New("неверное значение количества дней (должно быть от 1 до 400)")
		}

		// Обрабатываем добавление дней корректно
		nextDate := startDate.AddDate(0, 0, days)
		for !nextDate.After(now) {
			nextDate = nextDate.AddDate(0, 0, days)
		}

		// Проверка и корректировка 29 февраля
		if startDate.Month() == time.February && startDate.Day() == 29 {
			if !isLeapYear(startDate.Year()) {
				nextDate = time.Date(startDate.Year(), time.March, 1, 0, 0, 0, 0, nextDate.Location())
			} else if startDate.Year() != nextDate.Year() && nextDate.Month() == time.March && nextDate.Day() == 1 {
				nextDate = time.Date(nextDate.Year(), time.February, 28, 0, 0, 0, 0, nextDate.Location())
			}
		}

		return nextDate.Format("20060102"), nil
	}

	// Возвращаем ошибку для неподдерживаемых форматов
	return "", errors.New("неподдерживаемый формат повторения")
}

// isLeapYear определяет, является ли заданный год високосным.
func isLeapYear(year int) bool {
	return (year%4 == 0 && year%100 != 0) || (year%400 == 0)
}
