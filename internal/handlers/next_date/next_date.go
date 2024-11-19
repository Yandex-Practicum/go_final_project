package nextdate

import (
	"errors"
	"final_project/internal/common"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func NextDateHandler(w http.ResponseWriter, r *http.Request) {

	now, err := time.Parse(common.TimeFormat, r.FormValue("now"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	date := r.FormValue("date")
	repeat := r.FormValue("repeat")
	res, err := NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		fmt.Fprint(w, res)
	}

}

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("отсутствует правило повторения")
	}
	dateTime, err := time.Parse(common.TimeFormat, date)
	if err != nil {
		return "", errors.New("некорректная дата отсчета")
	}
	repeatParams := strings.Split(repeat, " ")
	switch repeatParams[0] {
	case "d":
		return daysRepeat(now, dateTime, repeatParams)
	case "y":
		return yearsRepeat(now, dateTime, repeatParams)
	default:
		return "", errors.New("некорректное правило повторения - ошибка в первом параметре")
	}
}

func daysRepeat(now time.Time, dateTime time.Time, repeatParams []string) (string, error) {
	if len(repeatParams) != 2 {
		return "", errors.New("некорректное правило повторения")
	}
	num, err := strconv.Atoi(repeatParams[1])
	if err != nil {
		return "", errors.New("некорректное правило повторения")
	}
	if num < 1 || num > 400 {
		return "", errors.New("некорректное правило повторения - число дней вне интервала (1-400)")
	}
	dateTime = dateTime.AddDate(0, 0, num)
	for dateTime.Before(now) {
		dateTime = dateTime.AddDate(0, 0, num)
	}
	return dateTime.Format(common.TimeFormat), nil
}

func yearsRepeat(now time.Time, dateTime time.Time, repeatParams []string) (string, error) {
	if len(repeatParams) != 1 {
		return "", errors.New("некорректное правило повторения")
	}
	dateTime = dateTime.AddDate(1, 0, 0)
	for dateTime.Before(now) {
		dateTime = dateTime.AddDate(1, 0, 0)
	}
	return dateTime.Format(common.TimeFormat), nil
}
