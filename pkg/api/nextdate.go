package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

// NextDate вычисляет следующую дату в зависимости от заданных параметров
func NextDate(nowParam, dateParam, repeatParam string) (string, error) {
	nowTime, err := parseTime(nowParam)
	if err != nil {
		return "", err
	}

	startDate, err := parseTime(dateParam)
	if err != nil {
		return "", err
	}

	repeatType, daysCount, err := parseRepeat(repeatParam)
	if err != nil {
		return "", err
	}

	var nextDate time.Time
	if repeatType == "y" {
		nextDate = getNextYearDate(nowTime, startDate)
	} else if repeatType == "d" {
		nextDate = getNextDayDate(nowTime, startDate, daysCount)
	} else {
		return "", fmt.Errorf("Неизвестный тип повторения")
	}

	return nextDate.Format(dateFormat), nil
}

// Обработчик для получения следующей даты
func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowParam := r.FormValue("now")
	dateParam := r.FormValue("date")
	repeatParam := r.FormValue("repeat")

	nextDate, err := NextDate(nowParam, dateParam, repeatParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nextDate))
}

// Помощные функции для парсинга
func parseTime(timeStr string) (time.Time, error) {
	if timeStr == "" {
		timeStr = time.Now().Format(dateFormat)
	}
	return time.Parse(dateFormat, timeStr)
}

func parseRepeat(repeat string) (string, int, error) {
	repeatParts := strings.Fields(repeat)
	if len(repeatParts) == 0 {
		return "", 0, fmt.Errorf("Параметр 'repeat' не указан")
	}

	if repeatParts[0] == "y" {
		return "y", 0, nil
	} else if repeatParts[0] == "d" {
		if len(repeatParts) != 2 {
			return "", 0, fmt.Errorf("Неправильный формат для 'd' (например, 'd 7')")
		}
		daysCount, err := strconv.Atoi(repeatParts[1])
		if err != nil || daysCount < 1 || daysCount > 400 {
			return "", 0, fmt.Errorf("Количество дней должно быть от 1 до 400")
		}
		return "d", daysCount, nil
	}
	return "", 0, fmt.Errorf("Неизвестный тип повторения")
}

func getNextYearDate(nowTime, startDate time.Time) time.Time {
	for {
		startDate = startDate.AddDate(1, 0, 0)
		if startDate.After(nowTime) {
			return startDate
		}
	}
}

func getNextDayDate(nowTime, startDate time.Time, daysCount int) time.Time {
	for {
		startDate = startDate.AddDate(0, 0, daysCount)
		if startDate.After(nowTime) {
			return startDate
		}
	}
}
