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
		for !AfterNow(date, now) {
			date = date.AddDate(0, 0, interval)
		}

	case "y":
		if len(parts) != 1 {
			return "", errors.New("неверный формат правила y")
		}
		for !AfterNow(date, now) {
			date = date.AddDate(1, 0, 0)
		}

	// можно оставить остальные правила w и m или вернуть ошибку, если не реализованы
	default:
		return "", fmt.Errorf("неподдерживаемый формат repeat: %s", repeat)
	}

	return date.Format("20060102"), nil
}
