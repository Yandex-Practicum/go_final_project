package api

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, dstartStr string, repeat string) (string, error) {

	var date time.Time

	//Нет повторения - ошибка
	if strings.TrimSpace(repeat) == "" {
		return "", errors.New("there is no repetition rule")
	}

	//Парсим
	dstart, err := time.Parse(Layout, dstartStr)
	if err != nil {
		return "", errors.New("incorrect start date")
	}

	//Правило повторения по частям
	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", errors.New("incorrect repetition rule")
	}

	ruleType := parts[0]

	//Обработка правила повторения (день, год)
	switch ruleType {
	case "d":
		if len(parts) != 2 {
			return "", errors.New("few parameters for d")
		}
		interval, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", errors.New("wrong format: " + parts[1])
		}
		if interval < 1 || interval > 400 {
			return "", errors.New("only number from 1 to 400 is supported")
		}

		date = dstart
		for {
			date = date.AddDate(0, 0, interval)
			if date.After(now) {
				return date.Format(Layout), nil
			}
		}

	case "y":
		if len(parts) != 1 {
			return "", errors.New("few parameters for y")
		}
		date = dstart
		for {
			date = date.AddDate(1, 0, 0)
			if date.After(now) {
				return date.Format(Layout), nil
			}
		}
	default:
		return "", fmt.Errorf("rule %s is not supported (yet)", ruleType)
	}

}
