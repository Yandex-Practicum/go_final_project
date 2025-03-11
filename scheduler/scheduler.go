package scheduler

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const DateLayout = "20060102"

// NextDate вычисляет следующую дату выполнения задачи, основываясь на базовых правилах:
// "d <число>" — прибавление указанного количества дней;
// "y" — ежегодное выполнение.
// Если правило не поддерживается или формат некорректен, возвращается ошибка.
func NextDate(now time.Time, dateStr, repeat string) (string, error) {
	baseDate, err := time.Parse(DateLayout, dateStr)
	if err != nil {
		return "", err
	}
	// Если базовая дата больше текущего времени, возвращаем её.
	if now.Before(baseDate) {
		return baseDate.Format(DateLayout), nil
	}

	repeat = strings.TrimSpace(repeat)
	if repeat == "" {
		return "", errors.New("Правило повторения не указано")
	}

	// Обработка правила "d <число>"
	if strings.HasPrefix(repeat, "d ") {
		parts := strings.Split(repeat, " ")
		if len(parts) != 2 {
			return "", errors.New("Неверный формат правила d")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days < 1 || days > 400 {
			return "", errors.New("Неверное количество дней")
		}
		nextDate := baseDate
		for !nextDate.After(now) {
			nextDate = nextDate.AddDate(0, 0, days)
		}
		return nextDate.Format(DateLayout), nil
	}

	// Обработка правила "y" – ежегодное выполнение.
	if repeat == "y" {
		nextDate := baseDate
		for !nextDate.After(now) {
			nextDate = nextDate.AddDate(1, 0, 0)
		}
		return nextDate.Format(DateLayout), nil
	}

	return "", errors.New("Неподдерживаемый формат правила повторения")
}
