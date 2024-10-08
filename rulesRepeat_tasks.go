package main

import (
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Обработчик для /api/nextdate
func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")
	now, err := time.Parse(dateFormat, nowStr)
	if err != nil {
		http.Error(w, "Неверный формат даты", http.StatusBadRequest)
		return
	}
	nextDate, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprint(w, nextDate)
}

// Вычисляем следующую дату для задачи на основе правил
func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("пусто")
	}
	tDate, err := time.Parse(dateFormat, date)
	if err != nil {
		return "", fmt.Errorf("неверный формат даты: %v", err)
	}
	startDate := now
	if tDate.After(now) {
		startDate = tDate
	}

	switch {
	case strings.HasPrefix(repeat, "d "): // ежедневное повторение "d"
		daysStr := strings.TrimSpace(repeat[2:])
		days, err := strconv.Atoi(daysStr)
		if err != nil || days <= 0 || days > 400 {
			return "", fmt.Errorf("неверноый формат: %v", err)
		}
		tDate = tDate.AddDate(0, 0, days)
		for tDate.Before(now) {
			tDate = tDate.AddDate(0, 0, days)
		}
		return tDate.Format(dateFormat), nil

	case repeat == "y": // ежегодное повторение "y"
		for !tDate.After(startDate) {
			year := tDate.Year() + 1
			month := tDate.Month()
			day := tDate.Day()

			// високосный год
			if month == time.February && day == 29 && !isLeapYear(year) {
				tDate = time.Date(year, time.March, 1, 0, 0, 0, 0, tDate.Location())
			} else {
				tDate = time.Date(year, month, day, 0, 0, 0, 0, tDate.Location())
			}
		}
		return tDate.Format(dateFormat), nil

	case strings.HasPrefix(repeat, "w "): // Еженедельное повторение
		daysStr := strings.TrimSpace(repeat[2:])
		days := strings.Split(daysStr, ",")
		if len(days) == 0 {
			return "", fmt.Errorf("дни не указаны")
		}

		var daysOfWeek []int
		for _, dayStr := range days {
			day, err := strconv.Atoi(dayStr)
			if err != nil || day < 1 || day > 7 {
				return "", fmt.Errorf("неверный день '%s'", dayStr)
			}
			if day == 7 {
				day = 0
			}
			daysOfWeek = append(daysOfWeek, day)
		}
		sort.Ints(daysOfWeek)

		startDate = tDate
		if now.After(tDate) {
			startDate = now
		}

		initialDate := tDate

		for !containsInt(daysOfWeek, int(startDate.Weekday())) || !(startDate.YearDay() > initialDate.YearDay()) {
			startDate = startDate.AddDate(0, 0, 1)
		}
		return startDate.Format(dateFormat), nil

	case strings.HasPrefix(repeat, "m "): // ежемесячное повторения
		parts := strings.Split(strings.TrimSpace(repeat[2:]), " ")

		if len(parts) == 0 {
			return "", fmt.Errorf("дни не указаны")
		}

		dayParts := strings.Split(parts[0], ",")
		var daysOfMonth []int
		for _, dayStr := range dayParts {
			day, err := strconv.Atoi(dayStr)
			if err != nil || day == 0 || day < -2 || day > 31 {
				return "", fmt.Errorf("неверный день '%s'", dayStr)
			}
			daysOfMonth = append(daysOfMonth, day)
		}

		var months []int
		if len(parts) > 1 {
			monthParts := strings.Split(parts[1], ",")
			for _, monthStr := range monthParts {
				month, err := strconv.Atoi(monthStr)
				if err != nil || month < 1 || month > 12 {
					return "", fmt.Errorf("неверный месяц '%s'", monthStr)
				}
				months = append(months, month)
			}
		}

		sort.Ints(daysOfMonth)
		sort.Ints(months)

		for {
			curYear, curMonth := tDate.Year(), tDate.Month()

			if len(months) > 0 && !containsInt(months, int(curMonth)) {
				nextMonth := findNextMonth(int(curMonth), months)
				if nextMonth <= int(curMonth) {
					tDate = time.Date(curYear+1, time.Month(nextMonth), 1, 0, 0, 0, 0, tDate.Location())
				} else {
					tDate = time.Date(curYear, time.Month(nextMonth), 1, 0, 0, 0, 0, tDate.Location())
				}
				continue
			}

			// Находим ближайшую допустимую дату в текущем месяце
			var nextValidDate time.Time
			for _, day := range daysOfMonth {
				var cDate time.Time
				lastDay := lastDayOfMonth(curYear, curMonth)

				if day > 0 {
					if day <= lastDay {
						cDate = time.Date(curYear, curMonth, day, 0, 0, 0, 0, tDate.Location())
					} else {
						continue
					}
				} else {
					if -day <= lastDay {
						cDate = time.Date(curYear, curMonth, lastDay+day+1, 0, 0, 0, 0, tDate.Location())
					} else {
						continue
					}
				}
				if cDate.After(startDate) && (nextValidDate.IsZero() || cDate.Before(nextValidDate)) {
					nextValidDate = cDate
				}
			}
			if !nextValidDate.IsZero() {
				return nextValidDate.Format(dateFormat), nil
			}
			tDate = tDate.AddDate(0, 1, 0)
			tDate = time.Date(tDate.Year(), tDate.Month(), 1, 0, 0, 0, 0, tDate.Location())
		}
	default:
		return "", fmt.Errorf("неверный формат: '%s'", repeat)
	}
}
