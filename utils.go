package main

import (
	"fmt"
	"time"
)

func NextDate(currentDate string, rule string) (string, error) {
	// Парсим текущую дату
	currentTime, err := time.Parse("2006-01-02", currentDate)
	if err != nil {
		return "", fmt.Errorf("не удалось разобрать текущую дату: %w", err)
	}

	// Рассчитываем следующую дату в зависимости от правила
	var nextDate time.Time

	switch rule {
	case "daily":
		nextDate = currentTime.Add(24 * time.Hour) // Добавить один день
	case "weekly":
		nextDate = currentTime.Add(7 * 24 * time.Hour) // Добавить одну неделю
	case "monthly":
		nextDate = currentTime.AddDate(0, 1, 0) // Добавить один месяц
	default:
		return "", fmt.Errorf("неподдерживаемое правило: %s", rule)
	}

	// Возвращаем следующую дату
	return nextDate.Format("2006-01-02"), nil
}
