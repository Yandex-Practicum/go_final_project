package main

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if len(repeat) == 0 {
		return "", errors.New("repeat is empty string")
	}

	dayMatched, _ := regexp.MatchString(`d \d{1,3}`, repeat)
	yearMatched, _ := regexp.MatchString(`y`, repeat)
	weekMatched, _ := regexp.MatchString(`w [1-7]+(,[1-7])*`, repeat)
	monthMatched, _ := regexp.MatchString(`m (\b(0?[1-9]|[1-2][0-9]|3[0-1]|-1|-2)\b|-1|-2)+(,\b(0?[1-9]|[1-2][0-9]|3[0-1])\b|,-1|,-2)* *(\b(0?[1-9]|1[0-2])\b)*(,\b(0?[1-9]|1[0-2])\b)*`, repeat)

	if dayMatched {
		days, err := strconv.Atoi(strings.TrimPrefix(repeat, "d "))
		if err != nil {
			return "", err
		}

		if days > 400 {
			return "", errors.New("maximum days count must be 400")
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
		// Извлечение дней месяца из строки repeat
		dayMonthStr := strings.Split(strings.TrimSpace(strings.TrimPrefix(repeat, "m")), " ")

		var daysOfMonth []int

		// Обработка дней месяца
		for _, d := range dayMonthStr {
			if d == "-1" {
				// Последний день месяца
				daysOfMonth = append(daysOfMonth, -1)
			} else if d == "-2" {
				// Предпоследний день месяца
				daysOfMonth = append(daysOfMonth, -2)
			} else {
				// Обычный день месяца
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

		// Извлечение месяцев из строки repeat (если указано)
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
			// Если месяц не указан, считаем, что задача должна быть назначена на все месяцы
			for i := 1; i <= 12; i++ {
				months = append(months, i)
			}
		}

		// Поиск следующей даты, удовлетворяющей условиям
		for _, month := range months {
			for _, day := range daysOfMonth {
				// Создание времени для проверки
				nextDate := time.Date(now.Year(), time.Month(month), day, 0, 0, 0, 0, time.Local)
				// Если дата раньше текущей, добавляем год
				if nextDate.Before(now) || nextDate.Equal(now) {
					nextDate = nextDate.AddDate(0, 1, 0)
				}
				// Проверка, не является ли следующая дата последним или предпоследним днем месяца
				if day == -1 || day == -2 {
					lastDay := time.Date(nextDate.Year(), nextDate.Month()+1, 0, 0, 0, 0, 0, time.Local)
					if day == -2 {
						lastDay = lastDay.AddDate(0, 0, -1)
					}
					nextDate = lastDay
				}
				// Если дата удовлетворяет текущему месяцу и условиям, возвращаем её
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

	// Проверка наличия параметров
	if now == "" || date == "" || repeat == "" {
		http.Error(w, "Missing parameters", http.StatusBadRequest)
		return
	}

	// Проверка формата параметров
	_, err := time.Parse("20060102", now)
	if err != nil {
		http.Error(w, "Invalid 'now' parameter format", http.StatusBadRequest)
		return
	}

	_, err = time.Parse("20060102", date)
	if err != nil {
		http.Error(w, "Invalid 'date' parameter format", http.StatusBadRequest)
		return
	}

	// Обработка параметра repeat
	nextDate, err := NextDate(time.Now(), date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK) // Вернуть код статуса 200
	w.Write([]byte(nextDate))
}
