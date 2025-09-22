package api

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Константа для парсинга
const dateFormat = "20060102"

func NextDate(now time.Time, dstartStr string, repeat string) (string, error) {

	var date time.Time

	//Нет повторения - ошибка
	if strings.TrimSpace(repeat) == "" {
		return "", errors.New("Нет повторения")
	}

	//Парсим
	dstart, err := time.Parse(dateFormat, dstartStr)
	if err != nil {
		return "", errors.New("Неверная дата")
	}

	//Правило повторения по частям
	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", errors.New("Неверный формат повторения")
	}

	ruleType := parts[0]

	//Обработка правила повторения (день, год)
	switch ruleType {
	case "d":
		if len(parts) != 2 {
			return "", errors.New("Мало параметров для d")
		}
		interval, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", errors.New("Неправильный формат: " + parts[1])
		}
		if interval < 1 || interval > 400 {
			return "", errors.New("Поддерживается только число от 1 до 400")
		}

		date = dstart
		for {
			date = date.AddDate(0, 0, interval)
			if date.After(now) {
				return date.Format(dateFormat), nil
			}
		}

	case "y":
		if len(parts) != 1 {
			return "", errors.New("Неправильный формат для y")
		}
		date = dstart
		for {
			date = date.AddDate(1, 0, 0)
			if date.After(now) {
				return date.Format(dateFormat), nil
			}
		}
	default:
		return "", fmt.Errorf("Правило %s не поддерживается(пока что)", ruleType)
	}

}
