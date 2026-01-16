package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var format string = "20060102"

func nextDayHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	// Получение данных
	now := req.FormValue("now")
	date := req.FormValue("date")
	repeat := req.FormValue("repeat")

	// Парсим из str в time
	parsedNow, err := time.Parse(format, now)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Вызов функции NextDate
	nextDate, err := NextDate(parsedNow, date, repeat)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Ответ
	w.Write([]byte(nextDate))
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("Нет повторений для задачи")
	}

	date, err := time.Parse(format, dstart)
	if err != nil {
		return "", err
	}

	flagAndTime := strings.Split(repeat, " ")
	if len(flagAndTime) > 2 {
		return "", errors.New("Неподдерживаемый формат повторений")
	}

	// В зависимости от флага ищем следующую дату повтора
	switch {
	case flagAndTime[0] == "d" && len(flagAndTime) == 2:
		d, err := strconv.Atoi(flagAndTime[1])
		if err != nil {
			return "", err
		}
		if d > 400 {
			return "", errors.New("Количство дней превышает 400")
		}
		if d < 1 {
			return "", errors.New("Количество дней меньше 1")
		}
		for {
			date = date.AddDate(0, 0, d)
			if afterNow(date, now) {
				break
			}
		}
	case flagAndTime[0] == "y" && len(flagAndTime) == 1:
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
		}
	default:
		return "", errors.New("Неподдерживаемый формат повторений")
	}

	return date.Format(format), nil
}

func afterNow(date, now time.Time) bool {
	// хотел сделать так, но тест, где первая дата - сегодня не проходит return date.After(now)

	return date.Format("20060102") > now.Format("20060102")
}
