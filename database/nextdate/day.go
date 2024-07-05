package nextdate

import (
	"fmt"
	"strconv"
)

// calcD возвращает следующую дату по правилу repeat "d"
func calcD(code string) (string, error) {
	days, err := strconv.Atoi(code)
	if err != nil {
		return "", err
	}
	if days > 400 {
		return "", fmt.Errorf("слишком большой временной промежуток")
	}
	nextDateDT := startDate.AddDate(0, 0, days)
	for nextDateDT.Before(now) || nextDateDT.Equal(now) {
		nextDateDT = nextDateDT.AddDate(0, 0, days)
	}
	nextDate = nextDateDT.Format(dateFormat)
	return nextDate, nil
}
