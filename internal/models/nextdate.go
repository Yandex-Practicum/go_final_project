package models

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, dateStr string, repeat string) (string, error) {
	date, err := time.Parse("20060102", dateStr)
	if err != nil {
		return "", errors.New("неверный формат даты")
	}

	if repeat == "" {
		return "", errors.New("правило повтора не указано")
	}

	switch {
	case repeat == "y":
		for {
			date = date.AddDate(1, 0, 0)
			if date.After(now) {
				break
			}
		}
	case strings.HasPrefix(repeat, "d "):
		days, err := strconv.Atoi(strings.TrimPrefix(repeat, "d "))
		if err != nil || days < 1 || days > 400 {
			return "", errors.New("неверный 'd' формат повтора")
		}
		for {
			date = date.AddDate(0, 0, days)
			if date.After(now) {
				break
			}
		}
	case strings.HasPrefix(repeat, "w "):
		return "", errors.New("неподдерживаемый формат повтора")
	case strings.HasPrefix(repeat, "m "):
		return "", errors.New("неподдерживаемый формат повтора")
	default:
		return "", errors.New("неверный формат повтора")
	}

	return date.Format("20060102"), nil
}
