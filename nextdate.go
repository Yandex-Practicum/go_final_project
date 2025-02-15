package main

import (
	"errors"
	"fmt"
	"time"
)

// Функция для вычисления следующей даты
func NextDate(now time.Time, date string, repeat string) (string, error) {
	parsedDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", errors.New("неверный формат даты")
	}

	switch {
	case repeat == "y":
		// Ежегодное повторение
		return parsedDate.AddDate(1, 0, 0).Format("20060102"), nil
	case len(repeat) > 2 && repeat[:2] == "d ":
		// Повторение через N дней
		var days int
		_, err := fmt.Sscanf(repeat, "d %d", &days)
		if err != nil || days > 400 {
			return "", errors.New("неверный формат правила повторения")
		}
		return parsedDate.AddDate(0, 0, days).Format("20060102"), nil
	default:
		return "", errors.New("неподдерживаемый формат правила повторения")
	}
}