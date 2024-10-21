package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)


func ApiNextDate(w http.ResponseWriter, r *http.Request) {
	nowStr := r.URL.Query().Get("now")
	dateStr := r.URL.Query().Get("date")
	repeatStr := r.URL.Query().Get("repeat")

	now, err := time.Parse(formatDate, nowStr)
	if err != nil {
		http.Error(w, "incorrect 'now' date format", http.StatusBadRequest)
		return
	}

	nextDate, err := NextDate(now, dateStr, repeatStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "%s", nextDate)

}

func NextDate(now time.Time, date string, repeat string) (string, error) {
	var res string
	if repeat == "" {
		return "", fmt.Errorf("repeat is empty")
	}

	validDate, err := time.Parse(formatDate, date)
	if err != nil {
		return "", fmt.Errorf("incorrect date format")
	}

	arr := strings.Split(repeat, " ")
	rule := arr[0]

	if rule == "d" {

		if len(arr) < 2 {
			return "", fmt.Errorf("no days")
		}
		days, err := strconv.Atoi(arr[1])
		if err != nil {
			return "", fmt.Errorf("days cant be int")
		}
		if days < 0 || days > 400 {
			return "", fmt.Errorf("incorrect days")
		}
		res = addDay(now, validDate, days)

	} else if rule == "y" {

		res = addYear(now, validDate)

	} else if rule == "w" {

		if len(arr) < 2 {
			return "", fmt.Errorf("no days")
		}
		res, err = addWeek(now, validDate, arr[1])
		if err != nil {
			return "", err
		}

	} else if rule == "m" {

		if len(arr) < 2 {
			return "", fmt.Errorf("no days")
		}
		res, err = addMonth(now, validDate, arr)
		if err != nil {
			return "", err
		}

	} else {
		return "", fmt.Errorf("incorrect rule")
	}

	return res, nil

}

func addDay(now, validDate time.Time, days int) string {
	if validDate.Equal(now) {
		return now.Format(formatDate)
	}
	validDate = validDate.AddDate(0, 0, days)
	for validDate.Before(now) {
		validDate = validDate.AddDate(0, 0, days)
	}
	return validDate.Format(formatDate)
}

func addYear(now, validDate time.Time) string {
	validDate = validDate.AddDate(1, 0, 0)
	for validDate.Before(now) {
		validDate = validDate.AddDate(1, 0, 0)
	}
	return validDate.Format(formatDate)
}

func addWeek(now, validDate time.Time, days string) (string, error) {
	arr := strings.Split(days, ",")

	weekDay := make(map[int]bool)

	for _, day := range arr {
		dayInt, err := strconv.Atoi(day)
		if err != nil {
			return "", err
		}
		if dayInt < 1 || dayInt > 7 {
			return "", fmt.Errorf("incorrect day of the week")
		}
		if dayInt == 7 {
			dayInt = 0
		}
		weekDay[dayInt] = true
	}

	for {
		if weekDay[int(validDate.Weekday())] && now.Before(validDate) {
			break
		}
		validDate = validDate.AddDate(0, 0, 1)
	}

	return validDate.Format(formatDate), nil

}

func addMonth(now, validDate time.Time, repeat []string) (string, error) {
	arr := strings.Split(repeat[1], ",")
	months := []string{}
	if len(repeat) > 2 {
		months = strings.Split(repeat[2], ",")
	}

	dayMap := map[int]bool{}

	for _, day := range arr {
		dayInt, err := strconv.Atoi(day)
		if err != nil {
			return "", err
		}
		if dayInt < -2 || dayInt > 31 || dayInt == 0 {
			return "", fmt.Errorf("incorrect day of the month")
		}
		dayMap[dayInt] = true
	}

	monthMap := map[int]bool{}
	for _, month := range months {
		if month == "" {
			continue
		}
		monthInt, err := strconv.Atoi(month)
		if err != nil {
			return "", err
		}
		if monthInt < 1 || monthInt > 12 {
			return "", fmt.Errorf("invalid month")
		}
		monthMap[monthInt] = true
	}

	for {
		if len(monthMap) == 0 || monthMap[int(validDate.Month())] {
			lastDay := time.Date(validDate.Year(), validDate.Month()+1, 0, 0, 0, 0, 0, validDate.Location()).Day()
			secondLastDay := lastDay - 1

			day := validDate.Day()
			if day == lastDay && dayMap[-1] {
				day = -1
			}
			if day == secondLastDay && dayMap[-2] {
				day = -2
			}

			if dayMap[day] && now.Before(validDate) {
				break
			}
		}
		validDate = validDate.AddDate(0, 0, 1)
	}
	return validDate.Format(formatDate), nil
}
