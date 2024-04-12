package main

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func addDays(date time.Time, days int) time.Time {
	return date.AddDate(0, 0, days)
}

func addYears(date time.Time, years int) time.Time {
	return date.AddDate(years, 0, 0)
}

func getNextDay(now time.Time, date string, repeat string) (string, error) {
	days, err := strconv.Atoi(strings.TrimPrefix(repeat, "d "))
	if err != nil {
		return "", err
	}

	if days > 400 {
		return "", errors.New("number of days exceeds the maximum limit of 400")
	}

	parsedDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", err
	}

	newDate := addDays(parsedDate, days)

	for newDate.Before(now) {
		newDate = addDays(newDate, days)
	}

	return newDate.Format("20060102"), nil
}

func getNextYear(now time.Time, date string) (string, error) {
	parsedDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", err
	}

	newDate := addYears(parsedDate, 1)

	for newDate.Before(now) {
		newDate = addYears(newDate, 1)
	}

	return newDate.Format("20060102"), nil
}

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if strings.HasPrefix(repeat, "d") {
		return getNextDay(now, date, repeat)
	} else if strings.Contains(repeat, "y") {
		return getNextYear(now, date)
	}
	return "", errors.New("repeat wrong format")
}

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	now := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	if now == "" || date == "" || repeat == "" {
		http.Error(w, "Missing parameters", http.StatusBadRequest)
		return
	}

	// Отформатируйте now в time.Time
	parsedNow, err := time.Parse("20060102", now)
	if err != nil {
		http.Error(w, "Invalid 'now' parameter", http.StatusBadRequest)
		return
	}

	// Вызовите функцию NextDate, передавая только день и год
	nextDate, err := NextDate(parsedNow, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nextDate))
}
