package nextdate

import (
	"fmt"
	"slices"
	"strings"
	"time"
)

func calcM(code string) (string, error) {
	// Текущая дата
	currentMonth := int(now.Month())
	currentDay := now.Day()
	currentYear := now.Year()

	// Разделяем repeat на указания
	codes := strings.Split(code, " ")
	// Если через пробел указано больше 2, или меньше 1 указания, то repeat был указан неправильно
	if len(codes) > 2 || len(codes) < 1 {
		return "", fmt.Errorf("некорректный формат repeat")
	}
	// Дни в которые должно происходить повторение
	targetDays, err := listAtoi(strings.Split(codes[0], ","))
	if err != nil {
		return "", err
	}

	// Проверяем, что эти дни не превышают максимальное количество дней в месяце
	idx := slices.IndexFunc(targetDays, func(d int) bool { return d > 31 || d < -2 })
	if idx != -1 {
		return "", fmt.Errorf("некорректный формат repeat")
	}
	// Проверяем, есть ли указания по месяцам, если есть то сохраняем их в targetMonths
	var targetMonths []int
	var monthSpecified bool
	if len(codes) > 1 {
		targetMonths, err = listAtoi(strings.Split(codes[1], ","))
		if err != nil {
			return "", err
		}
		idx := slices.IndexFunc(targetMonths, func(d int) bool { return d > 12 || d < 1 })
		if idx != -1 {
			return "", fmt.Errorf("некорректный формат repeat")
		}
		slices.Sort(targetMonths)
		monthSpecified = true
	}

	// Смотрим, будет ли ближайшая подходящая дата в этом месяце или году
	isNextMonth := getIsNextMonth(now, startDate, targetDays)
	isNextYear := getIsNextYear(isNextMonth, now, startDate, targetMonths)

	// Готовим переменные, из которых будет собирать следующую дату
	var nextYear, nextMonth, nextDay int
	// Собираем дату
	switch {

	// В этом месяце
	case !isNextMonth && (slices.Contains(targetMonths, currentMonth) || !monthSpecified):
		nextYear = currentYear
		nextMonth = currentMonth

		targetDays = processTargetDays(nextYear, time.Month(nextMonth), targetDays)
		// Поправка на дату начала
		var idx int
		startDay := pickBigger(currentDay, startDate.Day())
		if currentDay == startDate.Day() || currentDay == startDay {
			idx = getClosestIdx(startDay, targetDays)
			if idx == -1 {
				return "", fmt.Errorf("ошибка в вычисление ближайшего дня")
			}
		} else {
			idx = getClosestIdx(startDay-1, targetDays)
			if idx == -1 {
				return "", fmt.Errorf("ошибка в вычисление ближайшего дня")
			}
		}
		nextDay = targetDays[idx]

	// В этом году, в другом месяце, есть определённый месяц
	case !isNextYear && isNextMonth && monthSpecified:
		nextYear = currentYear
		// Поправка на дату начала
		var startMonth int
		if startDate.Year() < currentYear {
			startMonth = currentMonth
		} else {
			startMonth = pickBigger(currentMonth, int(startDate.Month()))
			if currentMonth < int(startDate.Month()) { // tbf??
				startMonth--
			}
		}

		idx := getClosestIdx(startMonth, targetMonths)
		if idx == -1 {
			return "", fmt.Errorf("ошибка в вычисление ближайшего месяца !isnextyear")
		}
		nextMonth = targetMonths[idx]

		targetDays = processTargetDays(nextYear, time.Month(nextMonth), targetDays)
		nextDay = targetDays[0]

	// В этом году, в другом месяце, нет требований к месяцу
	case !isNextYear && isNextMonth && !monthSpecified:
		nextYear = currentYear
		nextMonth = currentMonth + 1
		if startDate.Year() == currentYear {
			nextMonth = pickBigger(nextMonth, int(startDate.Month()))
		}
		targetDays = processTargetDays(nextYear, time.Month(nextMonth), targetDays)
		// Если минимальная возможная дата не существует в текущем месяце, берём следующий месяц
		for daysInMonth(nextYear, time.Month(nextMonth)) < targetDays[0] {
			nextMonth++
		}

		if nextMonth == int(startDate.Month()) {
			if idx = getClosestIdx(startDate.Day(), targetDays); idx != -1 {
				nextDay = targetDays[idx]
			} else {
				nextMonth++
				nextDay = targetDays[0]
			}

		} else { // Если следующий месяц больше чем начальный месяц, значит нам не нужно переживать при выборе дня, он будет позже начала
			nextDay = targetDays[0]
		}

	// В следующем году
	case isNextYear:
		if startDate.Year() > currentYear {
			nextYear = startDate.Year()
			if monthSpecified {
				if targetMonths[0] > int(startDate.Month()) {
					nextMonth = targetMonths[0]
				} else {
					idx := getClosestIdx(int(startDate.Month()), targetMonths)
					if idx == -1 {
						return "", fmt.Errorf("ошибка в вычисление ближайшего месяца isnextyear")
					}
					nextMonth = targetMonths[idx]
				}
			} else {
				nextMonth = int(startDate.Month())
			}
			targetDays = processTargetDays(nextYear, time.Month(nextMonth), targetDays)
			idx := getClosestIdx(startDate.Day(), targetDays)
			if idx == -1 {
				return "", fmt.Errorf("ошибка в вычисление ближайшего дня isnextyear")
			}
			nextDay = targetDays[idx]
		} else {
			nextYear = currentYear + 1
			if monthSpecified {
				nextMonth = targetMonths[0]
			} else {
				nextMonth = 1
			}
			targetDays = processTargetDays(nextYear, time.Month(nextMonth), targetDays)
			nextDay = targetDays[0]

		}
	default:
		return "", fmt.Errorf("ошибка в case m")
	}

	nextDate = time.Date(nextYear, time.Month(nextMonth), nextDay, 0, 0, 0, 0, time.UTC).Format(dateFormat)
	return nextDate, nil
}

// daysInMonth считает количество дней в указанном месяце указанного года
func daysInMonth(year int, month time.Month) int {
	t := time.Date(year, month, 32, 0, 0, 0, 0, time.UTC)
	daysInMonth := 32 - t.Day()
	return daysInMonth
}

// getIsNextMonth возвращает true если NextDate не может находится в текущем месяце
func getIsNextMonth(now time.Time, startDate time.Time, dates []int) bool {
	if startDate.Month() > now.Month() || startDate.Year() > now.Year() {
		return true
	}
	pdates := processTargetDays(now.Year(), now.Month(), dates)
	// Если getClosestIdx возвращает -1, значит в этом месяце не будет подходящего дня
	var startDay int
	if now.Month() > startDate.Month() {
		startDay = now.Day()
	} else {
		startDay = pickBigger(now.Day(), startDate.Day())

	}
	idx := getClosestIdx(startDay, pdates)

	return idx == -1

}

// getIsNextYear возвращает true если NextDate не может находится в текущем году
func getIsNextYear(isNextMonth bool, now time.Time, startDate time.Time, months []int) bool {
	if !isNextMonth {
		return false
	}
	if isNextMonth && int(now.Month()) == 12 {
		return true
	}
	if startDate.Year() > now.Year() {
		return true
	}

	if len(months) > 0 {
		idx := getClosestIdx(int(now.Month()), months)
		return idx == -1
	}
	return false
}

// processTargetDays возвращает отредактированный, отсортированный по возрастанию слайс с датами.
// Заменяет отрицательные числа на конкретные даты, исходя из количества дней в указанном месяце.
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
