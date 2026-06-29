package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

func daysInterval(splitedRepeat []string) (int, error) {
	if len(splitedRepeat) < 2 {
		err := errors.New("next date: Нарушение формата правила 'd' : пропущен интервал")
		return 0, err
	}
	days, err := strconv.Atoi(splitedRepeat[1])
	if err != nil {
		err = fmt.Errorf("next date: Ошибка конвертации строки в число %q: %w", splitedRepeat[1], err)
		return 0, err
	}
	if days < 1 || days > 400 {
		err = errors.New("next date: Указан неверный интервал, требуется от 1 до 400")
		return 0, err
	}
	return days, nil
}

func weeksInterval(splitedRepeat []string) (map[time.Weekday]bool, error) {
	if len(splitedRepeat) < 2 {
		err := errors.New("next date: Нарушение формата правила 'w': пропущен интервал")
		return nil, err
	}
	weekDays := make(map[time.Weekday]bool)
	weekDay := strings.Split(splitedRepeat[1], ",")
	for _, day := range weekDay {
		numOfDay, err := strconv.Atoi(day)
		if err != nil {
			err = fmt.Errorf("next date: Ошибка конвертации строки в число %q: %w", day, err)
			return nil, err
		}
		if numOfDay < 1 || numOfDay > 7 {
			err = errors.New("next date: Указан неверный интервал,требуется от 1 до 7")
			return nil, err
		}
		if numOfDay == 7 {
			numOfDay = 0
		}
		weekDays[time.Weekday(numOfDay)] = true
	}
	return weekDays, nil
}

func monthsInterval(splitedRepeat []string) (map[int]bool, map[int]bool, error) {
	if len(splitedRepeat) < 2 {
		err := errors.New("next date: Нарушение формата правила 'm': пропущен интервал")
		return nil, nil, err
	}
	months := make(map[int]bool)
	dayOfMonth := make(map[int]bool)

	daysOfMoth := strings.Split(splitedRepeat[1], ",")
	for _, day := range daysOfMoth {
		numOfDay, err := strconv.Atoi(day)
		if err != nil {
			err = fmt.Errorf("next date: Ошибка конвертации строки в число %q: %w", day, err)
			return nil, nil, err
		}
		if numOfDay < -2 || numOfDay > 31 || numOfDay == 0 {
			err = errors.New("next date: Указан неверный интервал,требуется от -2 до 31, не включая 0")
			return nil, nil, err
		}
		dayOfMonth[numOfDay] = true
	}
	if len(splitedRepeat) > 2 {
		spletedMonths := strings.Split(splitedRepeat[2], ",")
		for _, month := range spletedMonths {
			numOfMonth, err := strconv.Atoi(month)
			if err != nil {
				err = fmt.Errorf("next date: шибка конвертации строки в число %q: %w", month, err)
				return nil, nil, err
			}
			if numOfMonth < 1 || numOfMonth > 12 {
				err = errors.New("next date: Указан неверный интервал,требуется от 1 до 12")
				return nil, nil, err
			}
			months[numOfMonth] = true
		}
	}
	return dayOfMonth, months, nil
}

func AfterNow(nextDate time.Time, now time.Time) bool {
	nextStr := nextDate.Format(DateFormat)
	nowStr := now.Format(DateFormat)

	return nextStr > nowStr
}
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	date, err := time.Parse(DateFormat, dstart)
	if err != nil {
		err = fmt.Errorf("next date: Ошибка конвертации даты %q: %w", dstart, err)
		return "", err
	}

	if repeat == "" {
		return "", nil
	}
	splitedRepeat := strings.Split(repeat, " ")
	rule := splitedRepeat[0]
	nextDate := date
	switch rule {
	case "y":
		for {
			nextDate = nextDate.AddDate(1, 0, 0)
			if AfterNow(nextDate, now) {
				break
			}
		}
	case "d":
		days, err := daysInterval(splitedRepeat)
		if err != nil {
			return "", err
		}
		for {
			nextDate = nextDate.AddDate(0, 0, days)
			if AfterNow(nextDate, now) {
				break
			}
		}
	case "w":
		weekdays, err := weeksInterval(splitedRepeat)
		if err != nil {
			return "", err
		}
		for {
			nextDate = nextDate.AddDate(0, 0, 1)
			if weekdays[nextDate.Weekday()] && AfterNow(nextDate, now) {
				break
			}
		}
	case "m":
		dayOfMonth, months, err := monthsInterval(splitedRepeat)
		if err != nil {
			return "", err
		}
		for {
			nextDate = nextDate.AddDate(0, 0, 1)
			endOfMonth := (nextDate.AddDate(0, 1, -nextDate.Day())).Day()
			if len(months) == 0 {
				dayMatches := dayOfMonth[nextDate.Day()] ||
					(dayOfMonth[-1] && nextDate.Day() == endOfMonth) ||
					(dayOfMonth[-2] && nextDate.Day() == endOfMonth-1)
				if dayMatches && AfterNow(nextDate, now) {
					break
				}
			} else {
				dayMatches := months[int(nextDate.Month())] && (dayOfMonth[nextDate.Day()] ||
					(dayOfMonth[-1] && nextDate.Day() == endOfMonth) ||
					(dayOfMonth[-2] && nextDate.Day() == endOfMonth-1))
				if dayMatches && AfterNow(nextDate, now) {
					break
				}
			}
		}

	default:
		err := fmt.Errorf("next date: Указано неизвестное правило: %s", rule)
		return "", err
	}
	return nextDate.Format(DateFormat), nil
}

func NextDayHandler(res http.ResponseWriter, req *http.Request) {
	nowValue := req.FormValue("now")
	dstart := req.FormValue("date")
	repeat := req.FormValue("repeat")
	var now time.Time
	var err error
	if nowValue == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(DateFormat, nowValue)
		if err != nil {
			http.Error(res, "next date: Ошибка формата даты в переменной 'now' ", http.StatusBadRequest)
			return
		}
	}
	nextDate, err := NextDate(now, dstart, repeat)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	res.Write([]byte(nextDate))
}
