package api

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("repeat rule is empty")
	}

	startDate, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("invalid start date format: %v", err)
	}

	parts := strings.Split(repeat, " ")
	if len(parts) == 0 {
		return "", fmt.Errorf("invalid repeat format")
	}

	current := startDate
	nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid 'd' format, expected 'd <days>'")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days < 1 || days > 400 {
			return "", fmt.Errorf("invalid days interval, must be 1-400")
		}
		// Вычисляем разницу в днях между startDate и now
		currentDate := time.Date(current.Year(), current.Month(), current.Day(), 0, 0, 0, 0, current.Location())
		nowDateTruncated := time.Date(nowDate.Year(), nowDate.Month(), nowDate.Day(), 0, 0, 0, 0, nowDate.Location())
		if currentDate.After(nowDateTruncated) {
			// Если startDate > now, делаем один шаг
			current = current.AddDate(0, 0, days)
		} else {
			// Вычисляем, сколько дней прошло с startDate до now
			// Сколько полных интервалов прошло
			// Переходим к последней дате перед now
			// Добавляем ещё один интервал, чтобы превысить now
			daysSinceStart := int(nowDateTruncated.Sub(currentDate).Hours() / 24)
			intervalsPassed := daysSinceStart / days
			current = current.AddDate(0, 0, intervalsPassed*days)
			current = current.AddDate(0, 0, days)
		}
		return current.Format(DateFormat), nil

	case "y":
		currentDate := time.Date(current.Year(), current.Month(), current.Day(), 0, 0, 0, 0, current.Location())
		if !currentDate.After(nowDate) {
			for {
				current = current.AddDate(1, 0, 0)
				if current.Month() == time.February && current.Day() == 29 {
					year := current.Year()
					if year%4 != 0 || (year%100 == 0 && year%400 != 0) {
						current = time.Date(year, time.March, 1, 0, 0, 0, 0, current.Location())
					}
				}
				currentDate = time.Date(current.Year(), current.Month(), current.Day(), 0, 0, 0, 0, current.Location())
				if currentDate.After(nowDate) {
					break
				}
			}
		} else {
			current = current.AddDate(1, 0, 0)
			if current.Month() == time.February && current.Day() == 29 {
				year := current.Year()
				if year%4 != 0 || (year%100 == 0 && year%400 != 0) {
					current = time.Date(year, time.March, 1, 0, 0, 0, 0, current.Location())
				}
			}
		}
		return current.Format(DateFormat), nil

	case "w":
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid 'w' format, expected 'w <days>'")
		}
		daysStr := strings.Split(parts[1], ",")
		var days []int
		for _, d := range daysStr {
			day, err := strconv.Atoi(d)
			if err != nil || day < 1 || day > 7 {
				return "", fmt.Errorf("invalid day in 'w', must be 1-7")
			}
			days = append(days, day)
		}
		for {
			current = current.AddDate(0, 0, 1)
			currentDate := time.Date(current.Year(), current.Month(), current.Day(), 0, 0, 0, 0, current.Location())
			weekday := int(current.Weekday())
			if weekday == 0 {
				weekday = 7
			}
			for _, d := range days {
				if weekday == d && currentDate.After(nowDate) {
					return current.Format(DateFormat), nil
				}
			}
		}

	case "m":
		if len(parts) < 2 || len(parts) > 3 {
			return "", fmt.Errorf("invalid 'm' format, expected 'm <days> [months]'")
		}
		daysStr := strings.Split(parts[1], ",")
		var days []int
		for _, d := range daysStr {
			if d == "-1" {
				days = append(days, -1)
			} else if d == "-2" {
				days = append(days, -2)
			} else {
				day, err := strconv.Atoi(d)
				if err != nil || day < 1 || day > 31 {
					return "", fmt.Errorf("invalid day in 'm', must be 1-31 or -1,-2")
				}
				days = append(days, day)
			}
		}
		var months []int
		if len(parts) == 3 {
			monthsStr := strings.Split(parts[2], ",")
			for _, m := range monthsStr {
				month, err := strconv.Atoi(m)
				if err != nil || month < 1 || month > 12 {
					return "", fmt.Errorf("invalid month in 'm', must be 1-12")
				}
				months = append(months, month)
			}
		}
		for {
			current = current.AddDate(0, 0, 1)
			currentDate := time.Date(current.Year(), current.Month(), current.Day(), 0, 0, 0, 0, current.Location())
			day := current.Day()
			month := int(current.Month())
			lastDay := time.Date(current.Year(), current.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
			for _, d := range days {
				var match bool
				switch d {
				case -1:
					match = day == lastDay
				case -2:
					match = day == lastDay-1 && day > 0
				default:
					match = day == d
				}
				if match {
					if len(months) == 0 || contains(months, month) {
						if currentDate.After(nowDate) {
							return current.Format(DateFormat), nil
						}
					}
				}
			}
		}

	default:
		return "", fmt.Errorf("unsupported repeat rule: %s", parts[0])
	}
}

func contains(slice []int, item int) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

/*
func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	if nowStr == "" {
		now = time.Now()
	} else {
		var err error
		now, err = time.Parse(DateFormat, nowStr)
		if err != nil {
			http.Error(w, fmt.Sprintf("invalid now date: %v", err), http.StatusBadRequest)
			return
		}
	}

	if dateStr == "" {
		http.Error(w, "date parameter is required", http.StatusBadRequest)
		return
	}

	nextDate, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, nextDate)
}
*/
