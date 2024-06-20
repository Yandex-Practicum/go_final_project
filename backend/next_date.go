package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"
)

// daysBetweenWD возвращает количество дней между днями недели, с учётом их цикличности, в формате int (понедельник 1 ближе к воскресенью 7, чем пятница 5)
func daysBetweenWD(from, to int) int {
	week := make(map[int]int)
	for i := range 7 {
		if i+1 < 7 {
			week[i+1] = i + 2
		} else {
			week[i+1] = i + 2 - 7
		}
	}
	daysCount := 0
	i := week[from]
	for {
		if i == to {
			daysCount++
			return daysCount
		} else {
			daysCount++
			i = week[i]
		}
	}
}

// closestWD возвращает ближайший к текущему дню недели, день недели из списка, в формате int, с учётом цикличности недель
func closestWD(now time.Time, targetDays []int) int {
	if len(targetDays) < 2 {
		return -1
	}
	closestDay := 8
	currentDay := int(now.Weekday())
	for i := range targetDays {
		if daysBetweenWD(currentDay, targetDays[i]) < daysBetweenWD(currentDay, closestDay) {
			closestDay = targetDays[i]
		}
	}
	if closestDay > 7 || closestDay < 1 {
		return -1
	}
	return closestDay

}

// listAtoi конвертирует слайс string в слайс int
func listAtoi(list []string) ([]int, error) {
	var resList []int
	for _, str := range list {
		num, err := strconv.Atoi(str)
		if err != nil {
			return []int{}, err
		}
		resList = append(resList, num)
	}
	return resList, nil
}

// daysInMonth считает количество дней в указанном месяце указанного года
func daysInMonth(year int, month time.Month) int {
	t := time.Date(year, month, 32, 0, 0, 0, 0, time.UTC)
	daysInMonth := 32 - t.Day()
	return daysInMonth
}

// isNextMonth возвращает true если следующая дата не попадает в текущий месяц
func isNextMonth(now time.Time, dates []int) bool {
	maxDay := daysInMonth(now.Year(), now.Month())
	isNextMonth := true
	today := now.Day()
	for _, date := range dates {
		switch {
		case date > 0:
			if date > today && date <= maxDay {
				isNextMonth = false
				return isNextMonth
			}
		case date < 0:
			if today < (maxDay + date) {
				isNextMonth = false
				return isNextMonth
			}

		}

	}
	return isNextMonth
}

// isNextYear возвращает true если следующая дата не попадает в текущий год
func isNextYear(now time.Time, months []int) bool { // Заменить на поиск через slices.Index
	isNextYear := true
	currentMonth := int(now.Month())
	for _, month := range months {
		if month > currentMonth {
			isNextYear := false
			return isNextYear
		}
	}
	return isNextYear
}

// processTargetDays возвращает отредактированный слайс с датами, подходящими для следующей даты.
// заменяет отрицательные числа на конкретные даты, исходя из количества дней в месяце.
func processTargetDays(year int, month time.Month, target []int) []int {
	var processedTD []int
	daysTotal := daysInMonth(year, month)
	for _, tday := range target {
		if tday < 0 {
			day := daysTotal + tday + 1
			processedTD = append(processedTD, day)
		} else {
			processedTD = append(processedTD, tday)
		}
	}
	slices.Sort(processedTD)
	return processedTD
}

// getClosestIdx возвращает индекс блжиайшей переменной в слайсе list, превышающей значение target
func getClosestIdx(target int, list []int) int {
	idx := slices.IndexFunc(list, func(d int) bool { return d > target })
	return idx
}

// NextDate возвращает дату и ошибку, исходя из правил указанных в repeat.
func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("пустая строка")
	}

	format := "20060102"

	_, err := time.Parse(format, date)
	if err != nil {
		return "", err
	}

	var nextDate string = ""

	repeat = strings.ToLower(repeat)
	prefix := string(repeat[0])
	code, _ := strings.CutPrefix(repeat, prefix)
	code = strings.TrimSpace(code)

	switch prefix {
	case "d":
		days, err := strconv.Atoi(code)
		if err != nil {
			return "", err
		}
		if days > 400 {
			return "", fmt.Errorf("слишком большой временной промежуток")
		}
		nextDate = now.AddDate(0, 0, days).Format(format)

	case "y":
		lastDate, err := time.Parse(format, date)
		if err != nil {
			return "", err
		}
		nextDate = lastDate.AddDate(1, 0, 0).Format(format)

	case "w":
		currentWD := int(now.Weekday())
		targetWDs, err := listAtoi(strings.Split(code, ","))
		if err != nil {
			return "", err
		}

		closestWD := closestWD(now, targetWDs)
		days := daysBetweenWD(currentWD, closestWD)
		if closestWD == -1 || days == -1 {
			return "", fmt.Errorf("ошибка вычисления дней недели")
		}
		nextDate = now.AddDate(0, 0, days).Format(format)

	case "m":
		currentMonth := int(now.Month())
		currentDay := int(now.Day())
		currentYear := now.Year()

		codes := strings.Split(code, " ")
		if len(codes) > 2 || len(codes) < 1 {
			return "", fmt.Errorf("некорректный формат repeat")
		}
		targetDays, err := listAtoi(strings.Split(codes[0], ","))
		if err != nil {
			return "", err
		}
		idx := slices.IndexFunc(targetDays, func(d int) bool { return d > 31 || d < -31 })
		if idx != -1 {
			return "", fmt.Errorf("некорректный формат repeat")
		}
		var targetMonths []int
		var monthSpecified bool
		if len(code) > 1 {
			targetMonths, err = listAtoi(strings.Split(codes[1], ","))
			if err != nil {
				return "", err
			}
			idx := slices.IndexFunc(targetMonths, func(d int) bool { return d > 12 })
			if idx != -1 {
				return "", fmt.Errorf("некорректный формат repeat")
			}
			slices.Sort(targetMonths)
			monthSpecified = true
		}

		isNextMonth := isNextMonth(now, targetDays)
		isNextYear := isNextYear(now, targetMonths)

		var nextYear, nextMonth, nextDay int
		switch {
		// В этом месяце
		case !isNextMonth && (slices.Contains(targetMonths, currentMonth) || !monthSpecified):
			nextYear = currentYear
			nextMonth = currentMonth

			targetDays = processTargetDays(nextYear, time.Month(nextMonth), targetDays)
			idx := getClosestIdx(currentDay, targetDays)
			nextDay = targetDays[idx]

		// В этом году, в другом месяце
		case !isNextYear && !slices.Contains(targetMonths, currentMonth):
			nextYear = currentYear

			idx := getClosestIdx(currentMonth, targetMonths)
			nextMonth = targetMonths[idx]

			targetDays = processTargetDays(nextYear, time.Month(nextMonth), targetDays)
			nextDay = targetDays[0]
		// В следующем году
		case isNextYear:
			nextYear = currentYear + 1
			if monthSpecified {
				nextMonth = targetMonths[0]
			} else {
				nextMonth = 1
			}
			targetDays = processTargetDays(nextYear, time.Month(nextMonth), targetDays)
			nextDay = targetDays[0]
		}

		nextDate = time.Date(nextYear, time.Month(nextMonth), nextDay, 0, 0, 0, 0, time.UTC).Format(format)
	default:
		return "", fmt.Errorf("некорректный формат repeat")
	}

	return nextDate, nil

}
