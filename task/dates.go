package task

import (
	"errors"
	"strconv"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

// Парсим дату 4ый шаг
func ParseDate(dateStr string) (time.Time, error) {
	// Попробуем распарсить дату, если она передана
	if dateStr != "" {
		return time.Parse(DateFormat, dateStr)
	}
	// Если дата не передана, используем сегодняшнюю
	return time.Now(), nil
}

// Вычисление слудующей даты 3ий шаг
func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		err := errors.New("не указано правило повторения")
		return "", err
	}

	dateForm, err := time.Parse(DateFormat, date)
	if err != nil {
		err := errors.New("указан неверный формат времени")
		return "", err
	}
	rules := strings.Split(repeat, " ")

	switch rules[0] {
	case "d":
		if len(rules) != 2 {
			err = errors.New("неподдерживаемый формат правила повторения")
			return "", err
		}
		days, err := strconv.Atoi(rules[1])
		if err != nil {
			err = errors.New("неподдерживаемый формат правила повторения")
			return "", err
		}

		if days > 400 || days < 1 {
			err = errors.New("недопустимое количество дней")
			return "", err
		}
		for {
			dateForm = dateForm.AddDate(0, 0, days)
			if dateForm.After(now.Local().Truncate(24 * time.Hour)) {
				break
			}
		}
		return dateForm.Format(DateFormat), nil

	case "y":
		for {
			dateForm = dateForm.AddDate(1, 0, 0)
			if dateForm.After(now) {
				break
			}
		}
		return dateForm.Format(DateFormat), nil

	default:
		err = errors.New("недопустимое правило повторения")
		return "", err
	}
}
