package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {

	date, err := time.Parse(dateFormat, dstart)

	if err != nil {
		return "", fmt.Errorf("некорректная исходная дата: %w", err)
	}

	parts := strings.Fields(repeat)

	if len(parts) == 0 {
		return "", fmt.Errorf("правило повторения не указано")
	}

	switch parts[0] {

	case "d":

		if len(parts) != 2 {
			return "", fmt.Errorf("неверный формат правила d")
		}

		days, err := strconv.Atoi(parts[1])
		if err != nil || days < 1 || days > 400 {
			return "", fmt.Errorf("неверный интервал дней")
		}

		for {
			date = date.AddDate(0, 0, days)

			if afterNow(date, now) {
				break
			}
		}

	case "w":
		if len(parts) != 2 {
			return "", fmt.Errorf("неверный формат правила w")
		}

		weekDays := strings.Split(parts[1], ",")

		validDays := make(map[int]bool)

		for _, day := range weekDays {
			n, err := strconv.Atoi(day)
			if err != nil || n < 1 || n > 7 {
				return "", fmt.Errorf("неверный день недели: %s", day)
			}
			validDays[n] = true
		}

		for {
			date = date.AddDate(0, 0, 1)

			weekday := int(date.Weekday())
			if weekday == 0 {
				weekday = 7
			}

			if afterNow(date, now) && validDays[weekday] {
				break
			}
		}

	case "m":
		if len(parts) < 2 || len(parts) > 3 {
			return "", fmt.Errorf("неверный формат правила m")
		}

		monthDays := strings.Split(parts[1], ",")
		validDays := make(map[int]bool)

		for _, day := range monthDays {
			n, err := strconv.Atoi(day)
			if err != nil || (n < 1 || n > 31) && n != -1 && n != -2 {
				return "", fmt.Errorf("неверный день месяца: %s", day)
			}

			validDays[n] = true
		}

		validMonths := make(map[int]bool)

		if len(parts) == 3 {
			months := strings.Split(parts[2], ",")

			for _, month := range months {
				n, err := strconv.Atoi(month)
				if err != nil || n < 1 || n > 12 {
					return "", fmt.Errorf("неверный месяц: %s", month)
				}

				validMonths[n] = true
			}
		}

		for {
			date = date.AddDate(0, 0, 1)

			if len(parts) == 3 && !validMonths[int(date.Month())] {
				continue
			}

			lastDay := time.Date(
				date.Year(),
				date.Month()+1,
				0,
				0, 0, 0, 0,
				time.UTC,
			).Day()

			day := date.Day()

			if validDays[day] ||
				(validDays[-1] && day == lastDay) ||
				(validDays[-2] && day == lastDay-1) {

				if afterNow(date, now) {
					break
				}
			}
		}

	case "y":
		if len(parts) != 1 {
			return "", fmt.Errorf("неверный формат правила y")
		}

		for {
			date = date.AddDate(1, 0, 0)

			if afterNow(date, now) {
				break
			}
		}

	default:
		return "", fmt.Errorf("неподдерживаемое правило повторения: %s", parts[0])
	}

	return date.Format(dateFormat), nil
}

func afterNow(date time.Time, now time.Time) bool {
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	return date.After(now)
}

func nextDayHandler(w http.ResponseWriter, r *http.Request) {

	nowStr := r.FormValue("now")

	var now time.Time
	var err error

	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(dateFormat, nowStr)
		if err != nil {
			http.Error(w, "Некорректная дата now", http.StatusBadRequest)
			return
		}
	}

	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	nextDate, err := NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, nextDate)
}
