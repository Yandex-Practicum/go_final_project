package db

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("пустое правило повторения")
	}

	t, err := time.Parse(DateFormat, date)
	if err != nil {
		return "", fmt.Errorf("неверный формат даты: %s", err)
	}

	switch repeat[0] {
	case 'd':
		return nextDateByDays(now, t, repeat)
	case 'y':
		return nextDateByYear(now, t)
	case 'w':
		return nextDateByWeekday(now, t, repeat)
	case 'm':
		return nextDateByMonth(now, t, repeat)
	default:
		return "", fmt.Errorf("неизвестное правило повторения: %s", repeat)
	}
}

func nextDateByDays(now time.Time, t time.Time, repeat string) (string, error) {
	days, err := strconv.Atoi(strings.TrimSpace(repeat[1:]))
	if err != nil {
		return "", fmt.Errorf("неверный формат ежедневного повторения: %s", err)
	}
	if days < 1 || days > 400 {
		return "", fmt.Errorf("неверный формат ежедневного повторения: %d", days)
	}

	for {
		t = t.AddDate(0, 0, days)
		if t.After(now) {
			return t.Format(DateFormat), nil
		}
	}
}

func nextDateByYear(now time.Time, t time.Time) (string, error) {
	for {
		year := now.Year()
		if t.Year() < year {
			t = time.Date(year, t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
		} else {
			t = time.Date(t.Year()+1, t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
		}
		if t.After(now) {
			return t.Format(DateFormat), nil
		}
	}
}

func nextDateByWeekday(now time.Time, t time.Time, repeat string) (string, error) {
	parts := strings.Split(repeat, " ")
	if len(parts) != 2 {
		return "", fmt.Errorf("неверный формат повторения для правила 'w'")
	}

	weekdays := make([]int, 0)
	for _, day := range strings.Split(parts[1], ",") {
		dayInt, err := strconv.Atoi(day)
		if err != nil || dayInt < 1 || dayInt > 7 {
			return "", fmt.Errorf("неверный день недели в правиле 'w': %s", day)
		}
		weekdays = append(weekdays, dayInt)
	}

	nextDate := t
	for {
		weekday := int(nextDate.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		found := false
		for _, wd := range weekdays {
			if wd == weekday {
				found = true
				break
			}
		}
		if found {
			if nextDate.After(now) {
				return nextDate.Format(DateFormat), nil
			}
			nextDate = nextDate.AddDate(0, 0, 7)
		} else {
			nextDate = nextDate.AddDate(0, 0, 1)
		}
	}
}

func nextDateByMonth(now time.Time, t time.Time, repeat string) (string, error) {
	parts := strings.Split(repeat, " ")
	if len(parts) < 2 {
		return "", fmt.Errorf("неверный формат повторения для правила 'm'")
	}

	days := make([]int, 0)
	for _, day := range strings.Split(parts[1], ",") {
		dayInt, err := strconv.Atoi(day)
		if err != nil || (dayInt < 1 && dayInt != -1 && dayInt != -2) || dayInt > 31 {
			return "", fmt.Errorf("неверный день в правиле 'm': %s", day)
		}
		days = append(days, dayInt)
	}

	months := make([]int, 0)
	if len(parts) == 3 {
		for _, month := range strings.Split(parts[2], ",") {
			monthInt, err := strconv.Atoi(month)
			if err != nil || monthInt < 1 || monthInt > 12 {
				return "", fmt.Errorf("неверный месяц в правиле 'm': %s", month)
			}
			months = append(months, monthInt)
		}
	} else {
		for i := 1; i <= 12; i++ {
			months = append(months, i)
		}
	}

	nextDate := t
	for {
		day := nextDate.Day()
		month := int(nextDate.Month())
		lastDay := lastDayInMonth(nextDate)
		matchDay := false
		for _, d := range days {
			if d == day {
				matchDay = true
				break
			} else if d == -1 {
				if day == lastDay {
					matchDay = true
					break
				}
			} else if d == -2 {
				if day == lastDay-1 {
					matchDay = true
					break
				}
			}
		}
		matchMonth := true
		if len(months) > 0 {
			matchMonth = false
			for _, m := range months {
				if m == month {
					matchMonth = true
					break
				}
			}
		}
		if matchDay && matchMonth {
			if nextDate.After(now) {
				return nextDate.Format(DateFormat), nil
			}
			nextDate = time.Date(nextDate.Year(), nextDate.Month()+1, 1, 0, 0, 0, 0, time.UTC)
		} else {
			nextDate = nextDate.AddDate(0, 0, 1)
		}
	}
}

func lastDayInMonth(t time.Time) int {
	switch t.Month() {
	case time.February:
		if isLeapYear(t.Year()) {
			return 29
		} else {
			return 28
		}
	case time.April, time.June, time.September, time.November:
		return 30
	default:
		return 31
	}
}

func isLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}
