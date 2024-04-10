package main

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Repeater interface {
	NextDate(time.Time, time.Time) (time.Time, error)
}

type dayrepeater struct {
	days int
}

func (r *dayrepeater) NextDate(now time.Time, date time.Time) (time.Time, error) {
	next := date.AddDate(0, 0, r.days)
	for next.Before(now) {
		next = next.AddDate(0, 0, r.days)
	}
	return next, nil
}

func NewDay(repeat string) (*dayrepeater, error) {
	days, err := strconv.Atoi(strings.TrimPrefix(repeat, "d "))
	if err != nil {
		return nil, err
	}

	if days > 400 {
		return nil, errors.New("maximum days count must be 400")
	}

	return &dayrepeater{days}, nil
}

type yearrepeater struct {
	years int
}

func (r *yearrepeater) NextDate(now time.Time, date time.Time) (time.Time, error) {
	next := date.AddDate(r.years, 0, 0)
	for next.Before(now) {
		next = next.AddDate(r.years, 0, 0)
	}
	return next, nil
}

func NewYear(repeat string) (*yearrepeater, error) {
	return &yearrepeater{1}, nil
}

type weekrepeater struct {
	weekdays []int
}

func (r *weekrepeater) NextDate(now time.Time, date time.Time) (time.Time, error) {
	weekday := int(date.Weekday())
	var next time.Time
	updated := false
	for _, v := range r.weekdays {
		if weekday < v {
			next = date.AddDate(0, 0, v-weekday)
			updated = true
			break
		}
	}
	if !updated {
		next = date.AddDate(0, 0, 7-weekday+r.weekdays[0])
	}
	for next.Before(now) || next == now {
		weekday = int(next.Weekday())

		if weekday == r.weekdays[0] {
			for _, v := range r.weekdays {
				if weekday < v {
					next = next.AddDate(0, 0, v-weekday)
					weekday = int(next.Weekday())
				}
			}
		} else {
			next = next.AddDate(0, 0, 7-weekday+r.weekdays[0])
		}
	}
	return next, nil
}

func NewWeek(repeat string) (*weekrepeater, error) {
	var weekdays []int
	for _, val := range strings.Split(strings.TrimPrefix(repeat, "w "), ",") {
		int, err := strconv.Atoi(val)
		if err != nil {
			return nil, err
		}
		weekdays = append(weekdays, int)
	}
	return &weekrepeater{weekdays}, nil
}

type monthrepeater struct {
	daysOfMonth []int
	months      []int
}

func (r *monthrepeater) NextDate(now time.Time, date time.Time) (time.Time, error) {
	for _, month := range r.months {
		for _, day := range r.daysOfMonth {
			nextDate := time.Date(now.Year(), time.Month(month), day, 0, 0, 0, 0, time.Local)

			if nextDate.Before(now) || nextDate.Equal(now) {
				nextDate = nextDate.AddDate(0, 1, 0)
			}

			if nextDate.Month() == time.Month(month) && nextDate.After(now) {
				return nextDate, nil
			}
		}
	}

	return time.Time{}, errors.New("no suitable date found")
}

func NewMonth(repeat string) (*monthrepeater, error) {
	dayMonthStr := strings.Split(strings.TrimSpace(strings.TrimPrefix(repeat, "m")), " ")
	var daysOfMonth []int
	for _, d := range dayMonthStr {
		if d == "-1" || d == "-2" {
			return nil, errors.New("day of month should not be combined with last day indicator")
		} else {
			day, err := strconv.Atoi(d)
			if err != nil {
				return nil, err
			}
			if day < 1 || day > 31 {
				return nil, errors.New("invalid day of month")
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
				return nil, err
			}
			if month < 1 || month > 12 {
				return nil, errors.New("invalid month")
			}
			months = append(months, month)
		}
	} else {
		for i := 1; i <= 12; i++ {
			months = append(months, i)
		}
	}

	return &monthrepeater{daysOfMonth, months}, nil
}

func NextDate(now time.Time, date time.Time, repeat string) (string, error) {
	var repeater Repeater
	var err error
	switch {
	case strings.HasPrefix(repeat, "d"):
		repeater, err = NewDay(repeat)
	case strings.HasPrefix(repeat, "y"):
		repeater, err = NewYear(repeat)
	case strings.HasPrefix(repeat, "w"):
		repeater, err = NewWeek(repeat)
	case strings.HasPrefix(repeat, "m"):
		repeater, err = NewMonth(repeat)
	default:
		err = errors.New("unknown repeat")
	}

	if err != nil {
		return "", err
	}

	next, err := repeater.NextDate(now, date)
	if err != nil {
		return "", err
	}

	return next.Format("20060102"), nil
}

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	// Проверка формата параметров
	now, err := time.Parse("20060102", r.FormValue("now"))
	if err != nil {
		http.Error(w, "Invalid 'now' parameter format", http.StatusBadRequest)
		return
	}

	date, err := time.Parse("20060102", r.FormValue("date"))
	if err != nil {
		http.Error(w, "Invalid 'date' parameter format", http.StatusBadRequest)
		return
	}

	repeat := r.FormValue("repeat")
	if len(repeat) == 0 {
		http.Error(w, "repeat is empty string", http.StatusBadRequest)
	}

	// Обработка параметра repeat
	nextDate, err := NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK) // Вернуть код статуса 200
	w.Write([]byte(nextDate))
}
