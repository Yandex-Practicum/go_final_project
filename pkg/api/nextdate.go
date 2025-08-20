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
	if repeat == "" {
		return "", fmt.Errorf("repeat не должен быть пустым")
	}

	date, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("некорректная дата dstart: %v", err)
	}

	parts := strings.Split(repeat, " ")
	switch parts[0] {

	case "d":
		if len(parts) < 2 {
			return "", fmt.Errorf("отсутствует параметр для правила d")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", fmt.Errorf("неверный параметр для d: %v", err)
		}
		if days <= 0 || days > 400 {
			return "", fmt.Errorf("число дней для d должно быть от 1 до 400")
		}
		for {
			date = date.AddDate(0, 0, days)
			if afterNow(date, now) {
				break
			}
		}
		return date.Format(dateFormat), nil

	case "y":
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
		}
		return date.Format(dateFormat), nil

	case "w":
		if len(parts) < 2 {
			return "", fmt.Errorf("отсутствует список дней недели")
		}
		weekStrs := strings.Split(parts[1], ",")
		var weekdays [8]bool // от 1 до 7
		for _, w := range weekStrs {
			dayNum, err := strconv.Atoi(w)
			if err != nil || dayNum < 1 || dayNum > 7 {
				return "", fmt.Errorf("некорректный день недели: %v", w)
			}
			weekdays[dayNum] = true
		}
		for {
			date = date.AddDate(0, 0, 1)
			weekday := int(date.Weekday())
			if weekday == 0 {
				weekday = 7 // воскресенье
			}
			if weekdays[weekday] && date.After(now) {
				break
			}
		}
		return date.Format(dateFormat), nil

	case "m":
		if len(parts) < 2 {
			return "", fmt.Errorf("отсутствует список дней месяца")
		}
		daysStr := strings.Split(parts[1], ",")
		var dayFlags [32]bool
		hasMinus1 := false
		hasMinus2 := false

		for _, d := range daysStr {
			dayNum, err := strconv.Atoi(d)
			if err != nil {
				return "", fmt.Errorf("некорректный день месяца: %v", err)
			}
			if dayNum == -1 {
				hasMinus1 = true
			} else if dayNum == -2 {
				hasMinus2 = true
			} else if dayNum >= 1 && dayNum <= 31 {
				dayFlags[dayNum] = true
			} else {
				return "", fmt.Errorf("день месяца вне допустимого диапазона: %d", dayNum)
			}
		}

		var monthFlags [13]bool
		if len(parts) == 3 {
			monthStrs := strings.Split(parts[2], ",")
			for _, m := range monthStrs {
				monthNum, err := strconv.Atoi(m)
				if err != nil || monthNum < 1 || monthNum > 12 {
					return "", fmt.Errorf("месяц вне допустимого диапазона: %v", m)
				}
				monthFlags[monthNum] = true
			}
		} else {
			for i := 1; i <= 12; i++ {
				monthFlags[i] = true
			}
		}

		for {
			date = date.AddDate(0, 0, 1)
			day := date.Day()
			month := int(date.Month())
			year := date.Year()

			if !monthFlags[month] {
				continue
			}

			isMatch := false
			if dayFlags[day] {
				isMatch = true
			}
			if hasMinus1 {
				lastDay := time.Date(year, time.Month(month)+1, 0, 0, 0, 0, 0, time.UTC).Day()
				if day == lastDay {
					isMatch = true
				}
			}
			if hasMinus2 {
				lastDay := time.Date(year, time.Month(month)+1, 0, 0, 0, 0, 0, time.UTC).Day()
				if day == lastDay-1 {
					isMatch = true
				}
			}

			if isMatch && date.After(now) {
				break
			}
		}
		return date.Format(dateFormat), nil

	default:
		return "", fmt.Errorf("неподдерживаемый формат правила повторения: %s", parts[0])
	}
}

func nextDayHandler(w http.ResponseWriter, r *http.Request) {

	nowStr := r.FormValue("now")
	dstart := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	var err error
	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(dateFormat, nowStr)
		if err != nil {
			http.Error(w, fmt.Sprintf("неверный параметр now: %v", err), http.StatusBadRequest)
			return
		}
	}

	next, err := NextDate(now, dstart, repeat)
	if err != nil {
		http.Error(w, fmt.Sprintf("ошибка вычисления следующей даты: %v", err), http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, next)
}
func afterNow(date, now time.Time) bool {
	y1, m1, d1 := date.Date()
	y2, m2, d2 := now.Date()
	dateZero := time.Date(y1, m1, d1, 0, 0, 0, 0, time.UTC)
	nowZero := time.Date(y2, m2, d2, 0, 0, 0, 0, time.UTC)
	return !dateZero.Before(nowZero)
}
