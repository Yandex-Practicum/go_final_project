package main

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	dayMatched := strings.HasPrefix(repeat, "d")
	yearMatched := strings.Contains(repeat, "y")
	weekMatched := strings.HasPrefix(repeat, "w")
	monthMatched := strings.HasPrefix(repeat, "m")

	if dayMatched {
		days, err := strconv.Atoi(strings.TrimPrefix(repeat, "d "))
		if err != nil {
			return "", err
		}

		parsedDate, err := time.Parse("20060102", date)
		if err != nil {
			return "", err
		}

		newDate := parsedDate.AddDate(0, 0, days)

		for newDate.Before(now) {
			newDate = newDate.AddDate(0, 0, days)
		}

		return newDate.Format("20060102"), nil
	} else if yearMatched {
		parsedDate, err := time.Parse("20060102", date)
		if err != nil {
			return "", err
		}

		newDate := parsedDate.AddDate(1, 0, 0)

		for newDate.Before(now) {
			newDate = newDate.AddDate(1, 0, 0)
		}

		return newDate.Format("20060102"), nil
	} else if weekMatched {
		parsedDate, err := time.Parse("20060102", date)
		weekday := int(parsedDate.Weekday())
		if err != nil {
			return "", err
		}

		var newDate time.Time
		var weekdays []int

		for _, weekdayString := range strings.Split(strings.TrimPrefix(repeat, "w "), ",") {
			weekdayInt, _ := strconv.Atoi(weekdayString)
			weekdays = append(weekdays, weekdayInt)
		}

		updated := false
		for _, v := range weekdays {
			if weekday < v {
				newDate = parsedDate.AddDate(0, 0, v-weekday)
				updated = true
				break
			}
		}

		if !updated {
			newDate = parsedDate.AddDate(0, 0, 7-weekday+weekdays[0])
		}

		for newDate.Before(now) || newDate == now {
			weekday = int(newDate.Weekday())

			if weekday == weekdays[0] {
				for _, v := range weekdays {
					if weekday < v {
						newDate = newDate.AddDate(0, 0, v-weekday)
						weekday = int(newDate.Weekday())
					}
				}
			} else {
				newDate = newDate.AddDate(0, 0, 7-weekday+weekdays[0])
			}
		}

		return newDate.Format("20060102"), nil
	} else if monthMatched {
		dayMonthStr := strings.Split(strings.TrimSpace(strings.TrimPrefix(repeat, "m")), " ")

		var daysOfMonth []int

		for _, d := range dayMonthStr {
			if d == "-1" || d == "-2" {
				return "", errors.New("day of month should not be combined with last day indicator")
			} else {
				day, err := strconv.Atoi(d)
				if err != nil {
					return "", err
				}
				if day < 1 || day > 31 {
					return "", errors.New("invalid day of month")
				}
				daysOfMonth = append(daysOfMonth, day)
			}
		}

		monthStr := strings.TrimSpace(strings.SplitN(repeat, " ", 3)[2])

		var months []int

		if monthStr != "" {
			monthList := strings.Split(monthStr, ",")
			for _, m := range monthList {
				month, err := strconv.Atoi(m)
				if err != nil {
					return "", err
				}
				if month < 1 || month > 12 {
					return "", errors.New("invalid month")
				}
				months = append(months, month)
			}
		} else {
			for i := 1; i <= 12; i++ {
				months = append(months, i)
			}
		}

		for _, month := range months {
			for _, day := range daysOfMonth {
				nextDate := time.Date(now.Year(), time.Month(month), day, 0, 0, 0, 0, time.Local)

				if nextDate.Before(now) || nextDate.Equal(now) {
					nextDate = nextDate.AddDate(0, 1, 0)
				}

				if nextDate.Month() == time.Month(month) && nextDate.After(now) {
					return nextDate.Format("20060102"), nil
				}
			}
		}

		return "", errors.New("no suitable date found")
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

	nextDate, err := NextDate(time.Now(), date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nextDate))
}
