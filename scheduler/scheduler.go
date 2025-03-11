package scheduler

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const DateLayout = "20060102"

// addYearAdjusted добавляет год к дате с корректировкой для 29 февраля.
func addYearAdjusted(t time.Time) time.Time {
	res := t.AddDate(1, 0, 0)
	// Если исходная дата была 29 февраля, а результат стал 28 февраля, то устанавливаем 1 марта.
	if t.Month() == time.February && t.Day() == 29 {
		if res.Month() == time.February && res.Day() == 28 {
			res = time.Date(res.Year(), time.March, 1, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
		}
	}
	return res
}

func NextDate(now time.Time, dateStr, repeat string) (string, error) {
	baseDate, err := time.Parse(DateLayout, dateStr)
	if err != nil {
		return "", err
	}

	repeat = strings.TrimSpace(repeat)
	if repeat == "" {
		return "", errors.New("Правило повторения не указано")
	}

	// Правило "d <число>"
	if strings.HasPrefix(repeat, "d ") {
		parts := strings.Split(repeat, " ")
		if len(parts) != 2 {
			return "", errors.New("Неверный формат правила d")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days < 1 || days > 400 {
			return "", errors.New("Неверное количество дней")
		}
		// Всегда добавляем интервал хотя бы один раз.
		nextDate := baseDate.AddDate(0, 0, days)
		for !nextDate.After(now) {
			nextDate = nextDate.AddDate(0, 0, days)
		}
		return nextDate.Format(DateLayout), nil
	}

	// Правило "y" – ежегодное выполнение.
	if repeat == "y" {
		nextDate := addYearAdjusted(baseDate)
		for !nextDate.After(now) {
			nextDate = addYearAdjusted(nextDate)
		}
		return nextDate.Format(DateLayout), nil
	}

	return "", errors.New("Неподдерживаемый формат правила повторения")
}
