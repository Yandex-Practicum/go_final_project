package api

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, dstart string, repeat string) (string, error) {

	if repeat == "" {
		return "", fmt.Errorf("правило повторения не может быть пустым")
	}

	date, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("неверный формат даты: %s", dstart)
	}

	parts := strings.Fields(repeat)
	rule := parts[0]

	switch rule {
	case "d":

		if len(parts) != 2 {
			return "", fmt.Errorf("формат правила d неверный, ожидался: d <число>")
		}

		days, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", fmt.Errorf("неверное количество дней: %s", parts[1])
		}

		if days < 1 {
			return "", fmt.Errorf("количество дней должно быть больше 0, получено: %d", days)
		}

		if days > 400 {
			return "", fmt.Errorf("превышен максимально допустимый интервал: %d (максимум 400)", days)
		}

		nextDate := date
		for {
			nextDate = nextDate.AddDate(0, 0, days)
			if nextDate.After(now) {
				return nextDate.Format(DateFormat), nil
			}
		}

	case "y":

		if len(parts) != 1 {
			return "", fmt.Errorf("формат правила y неверный, лишние параметры")
		}

		nextDate := date
		for {
			nextDate = nextDate.AddDate(1, 0, 0)
			if nextDate.After(now) {
				return nextDate.Format(DateFormat), nil
			}
		}

	default:

		return "", fmt.Errorf("неподдерживаемый формат правила: %s (поддерживаются: d, y)", rule)
	}
}
