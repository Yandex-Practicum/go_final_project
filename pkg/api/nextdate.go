package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

func afterNow(date, now time.Time) bool {
	return date.After(now)
}

// parseNumbers парсит строку чисел, разделённых запятыми, в срез int
func parseNumbers(s string) ([]int, error) {
	if s == "" {
		return nil, errors.New("empty number list")
	}
	numStrs := strings.Split(s, ",")
	nums := make([]int, 0, len(numStrs))
	for _, numStr := range numStrs {
		num, err := strconv.Atoi(strings.TrimSpace(numStr))
		if err != nil {
			return nil, errors.New("invalid number format")
		}
		nums = append(nums, num)
	}
	return nums, nil
}

// isValidWeekDay проверяет, что дни недели в пределах 1–7
func isValidWeekDay(days []int) bool {
	for _, day := range days {
		if day < 1 || day > 7 {
			return false
		}
	}
	return true
}

// isValidMonthDay проверяет, что дни месяца в пределах 1–31 или -1, -2
func isValidMonthDay(days []int) bool {
	for _, day := range days {
		if day < -2 || day == 0 || day > 31 {
			return false
		}
	}
	return true
}

// isValidMonth проверяет, что месяцы в пределах 1–12
func isValidMonth(months []int) bool {
	for _, month := range months {
		if month < 1 || month > 12 {
			return false
		}
	}
	return true
}

// getLastDayOfMonth возвращает последний день месяца
func getLastDayOfMonth(t time.Time) int {
	return t.AddDate(0, 1, -t.Day()).Day()
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	// Парсим дату старта
	date, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", errors.New("invalid date format")
	}

	// Если repeat пустой, возвращаем ошибку
	if repeat == "" {
		return "", errors.New("repeat rule is empty")
	}

	// Разбираем repeat на части
	parts := strings.Split(repeat, " ")
	if len(parts) == 0 {
		return "", errors.New("repeat rule is empty")
	}

	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", errors.New("invalid d rule format")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days <= 0 || days > 400 {
			return "", errors.New("invalid number of days")
		}

		if !date.Before(now) {
			date = date.AddDate(0, 0, days)
			return date.Format(dateFormat), nil
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
		if len(parts) != 2 {
			return "", errors.New("invalid w rule format")
		}
		weekDays, err := parseNumbers(parts[1])
		if err != nil {
			return "", errors.New("invalid week days format")
		}
		if !isValidWeekDay(weekDays) {
			return "", errors.New("week days must be between 1 and 7")
		}
		// Продолжаем искать подходящий день недели
		for {
			date = date.AddDate(0, 0, 1)
			weekDay := int(date.Weekday())
			if weekDay == 0 {
				weekDay = 7 // Воскресенье
			}
			for _, wd := range weekDays {
				if wd == weekDay && afterNow(date, now) {
					return date.Format(dateFormat), nil
				}
			}
		}

	case "m":
		var monthDays, months []int
		if len(parts) < 2 || len(parts) > 3 {
			return "", errors.New("invalid m rule format")
		}
		monthDays, err = parseNumbers(parts[1])
		if err != nil {
			return "", errors.New("invalid month days format")
		}
		if !isValidMonthDay(monthDays) {
			return "", errors.New("month days must be between 1 and 31 or -1, -2")
		}
		if len(parts) == 3 {
			months, err = parseNumbers(parts[2])
			if err != nil {
				return "", errors.New("invalid months format")
			}
			if !isValidMonth(months) {
				return "", errors.New("months must be between 1 and 12")
			}
		} else {
			// Если месяцы не указаны, используем все
			months = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
		}
		// Начинаем с dstart
		for {
			date = date.AddDate(0, 0, 1)
			currentMonth := int(date.Month())
			currentDay := date.Day()
			for _, month := range months {
				if currentMonth != month {
					continue
				}
				for _, md := range monthDays {
					if md > 0 && md == currentDay {
						if afterNow(date, now) {
							return date.Format(dateFormat), nil
						}
					} else if md < 0 {
						lastDay := getLastDayOfMonth(date)
						if currentDay == lastDay+md+1 {
							if afterNow(date, now) {
								return date.Format(dateFormat), nil
							}
						}
					}
				}
			}
			// Защита от бесконечного цикла (например, 100 лет)
			if date.After(now.AddDate(100, 0, 0)) {
				return "", errors.New("no valid month day found within reasonable time")
			}
		}

	default:
		return "", errors.New("unsupported repeat rule")
	}
}

// nextDayHandler обрабатывает HTTP-запрос и вычисляет следующую дату по правилам
func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	// Если параметр now не передан, используем текущую дату
	if nowStr == "" {
		nowStr = time.Now().Format(dateFormat)
	}

	// Парсим дату now
	now, err := time.Parse(dateFormat, nowStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid now parameter: %v", err), http.StatusBadRequest)
		return
	}

	// Вызываем функцию для вычисления следующей даты
	nextDate, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, fmt.Sprintf("error: %v", err), http.StatusBadRequest)
		return
	}

	// Возвращаем результат
	w.Write([]byte(nextDate))
}
