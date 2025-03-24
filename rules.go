package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var timeFormat = "20060102"

// Формирует мапу допустимых значений дней в месяце
func daysInMonth(year int) map[int]int {
	daysAmout := make(map[int]int)
	for month := 1; month < 13; month++ {
		nextMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
		days := nextMonth.AddDate(0, 0, -1).Day()
		daysAmout[month] = days
	}
	return daysAmout
}

// ParseRepeater парсит строку на символ, слайс первой группы чисел, слайс второй группы чисел и ошибку
func ParseRepeater(repeat string) (string, []int, []int, error) {
	if repeat == "" {
		return "", nil, nil, fmt.Errorf("пустая строка")
	}

	// Обработка для месяцев
	if strings.HasPrefix(repeat, "m ") {
		symbol := "m"
		monthStr := strings.TrimPrefix(repeat, "m ")
		parts := strings.Split(monthStr, " ")

		if len(parts) < 1 {
			return "", nil, nil, fmt.Errorf("некорректный ввод для месяца")
		}

		// Обработка первой группы чисел (дни месяца)
		var daysOfMonth []int
		dayStr := strings.Split(parts[0], ",")
		for _, day := range dayStr {
			value, err := strconv.Atoi(day)
			if err != nil {
				return "", nil, nil, fmt.Errorf("некорректное число дня: %s", day)
			}
			daysOfMonth = append(daysOfMonth, value)
		}

		// Обработка второй группы чисел (месяцы)
		var months []int
		if len(parts) > 1 {
			monthStr := strings.Join(parts[1:], " ")
			monthNumbers := strings.Split(monthStr, ",")
			for _, month := range monthNumbers {
				value, err := strconv.Atoi(month)
				if err != nil {
					return "", nil, nil, fmt.Errorf("некорректное число месяца: %s", month)
				}
				months = append(months, value)
			}
		}

		return symbol, daysOfMonth, months, nil
	}

	// Парсинг repeat на символ и числа
	elements := strings.Split(repeat, " ")
	if len(elements) == 0 {
		return "", nil, nil, fmt.Errorf("некорректный ввод")
	}

	// Получения символа repeat
	symbol := elements[0]
	var days []int

	// Создания слайса из чисел repeat
	if len(elements) > 1 {
		nums := strings.Split(elements[len(elements)-1], ",")
		for _, num := range nums {
			day, err := strconv.Atoi(num)
			if err != nil {
				return "", nil, nil, err
			}
			days = append(days, day)
		}
	}

	return symbol, days, nil, nil
}

// NextDate возвращает слайс, содержащий даты задач в формате "20060102"
func NextDate(now time.Time, date string, repeat string) ([]string, error) {
	var (
		futureDate time.Time
		taskDays   []string
	)

	t, err := time.Parse(timeFormat, date)
	if err != nil {
		return nil, fmt.Errorf("ошибка при парсинге даты: %v", err)
	}

	symbol, days, months, err := ParseRepeater(repeat)
	if err != nil {
		return nil, err
	}

	switch symbol {
	case "d":
		if len(days) == 0 {
			return nil, fmt.Errorf("не указан интервал в днях")
		}

		if len(days) > 1 {
			return nil, fmt.Errorf("недопустимое количество дней")
		}

		if days[0] > 400 {
			return nil, fmt.Errorf("превышен максимально допустимый интервал дней")
		}

		futureDate = t.AddDate(0, 0, days[0])
		for futureDate.Before(now) {
			futureDate = futureDate.AddDate(0, 0, days[0])
		}

		taskDays = append(taskDays, futureDate.Format(timeFormat))
		return taskDays, nil

	case "w":
		if len(days) == 0 {
			return nil, fmt.Errorf("не указан день недели")
		}

		for _, day := range days {
			if day < 1 || day > 7 {
				return nil, fmt.Errorf("недопустимое значение дня недели %v", day)
			}
		}

		currentWeekDay := int(now.Weekday())
		if currentWeekDay == 0 {
			currentWeekDay = 7
		}

		for _, day := range days {
			dayDif := (day - currentWeekDay + 7) % 7
			if dayDif == 0 {
				dayDif = 7
			}
			futureDate = now.AddDate(0, 0, int(dayDif))

			taskDays = append(taskDays, futureDate.Format(timeFormat))
		}

	case "m":
		if len(days) == 0 {
			return nil, fmt.Errorf("не указан день месяца")
		}

		for _, day := range days {
			if day > 31 || day < -2 {
				return nil, fmt.Errorf("недопустимое значение числа месяца %v", day)
			}
		}

		currentMonthDay := t.Day()
		currentMonth := t.Month()
		currentYear := t.Year()

		// Получаем мапу допустимых значений
		daysAmount := daysInMonth(currentYear)

		// Если есть вторая группа чисел после "m"
		if months != nil {
			for _, month := range months {
				if month < 1 || month > 12 {
					return nil, fmt.Errorf("недопустимое значение месяца %v", month)
				}

				for _, day := range days {
					if day > 0 && day <= currentMonthDay { // Если день меньше или равен текущему дню, переносится на следующий месяц
						futureDate = time.Date(currentYear, time.Month(month)+1, day, 0, 0, 0, 0, now.Location())
					} else if day < 0 { // Если указано -1 или -2 переносится на последнее и предпоследнее число месяца, соотвественно
						futureDate = time.Date(currentYear, time.Month(month)+1, 1, 0, 0, 0, 0, now.Location()).AddDate(0, 0, day)
					} else { // Если день больше текущего дня, используется текущий месяц
						futureDate = time.Date(currentYear, time.Month(month), day, 0, 0, 0, 0, now.Location())
					}
				}
			}
		} else { // если только одна группа чисел после "m"
			for _, day := range days {
				if day > daysAmount[int(currentMonth)] {
					futureDate = time.Date(currentYear, t.Month()+1, day, 0, 0, 0, 0, now.Location())
				}
				if day > 0 && day <= currentMonthDay { // Если день меньше или равен текущему дню, переносится на следующий месяц
					futureDate = time.Date(currentYear, t.Month()+1, day, 0, 0, 0, 0, now.Location())
				} else if day < 0 { // Если указано -1 или -2 переносится на последнее и предпоследнее число месяца, соотвественно
					futureDate = time.Date(currentYear, t.Month()+1, 1, 0, 0, 0, 0, now.Location()).AddDate(0, 0, day)
				} else { // Если день больше текущего дня, используется текущий месяц
					futureDate = time.Date(currentYear, t.Month(), day, 0, 0, 0, 0, now.Location())
				}
				// for futureDate.Before(now) {
				// 	futureDate = futureDate.AddDate(0, 1, 0)
				// }
				taskDays = append(taskDays, futureDate.Format(timeFormat))
			}
		}

	case "y":
		if len(days) != 0 {
			return nil, fmt.Errorf("недопустимое значение для года")
		}

		futureDate = t.AddDate(1, 0, 0)
		for futureDate.Before(now) {
			futureDate = futureDate.AddDate(1, 0, 0)
		}

		taskDays = append(taskDays, futureDate.Format(timeFormat))
		return taskDays, nil

	default:
		return nil, fmt.Errorf("недопустимый символ")

	}

	return taskDays, nil
}
