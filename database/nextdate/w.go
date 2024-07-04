package nextdate

import (
	"fmt"
	"slices"
	"strings"
	"time"
)

// calcW возвращает следующую дату или ошибку, исходя из правила repeat "w"
func calcW(code string) (string, error) {
	var days int
	targetWDs, err := listAtoi(strings.Split(code, ","))
	if err != nil {
		return "", err
	}
	idx := slices.IndexFunc(targetWDs, func(wd int) bool { return wd > 7 || wd < 1 })
	if idx != -1 {
		return "", fmt.Errorf("некорректный формат w")
	}

	switch {
	case startDate.After(now) || startDate.Equal(now):
		startWD := int(startDate.Weekday())
		closestWD := closestWD(startDate, targetWDs)
		days = daysBetweenWD(startWD, closestWD)
		if closestWD == -1 || days == -1 {
			return "", fmt.Errorf("ошибка вычисления дней недели")
		}
	case startDate.Before(now):
		currentWD := int(now.Weekday())

		closestWD := closestWD(now, targetWDs)
		days = daysBetweenWD(currentWD, closestWD)
		if closestWD == -1 || days == -1 {
			return "", fmt.Errorf("ошибка вычисления дней недели")
		}
	}
	nextDate = now.AddDate(0, 0, days).Format(dateFormat)
	return nextDate, err

}

// daysBetweenWD возвращает количество дней между днями недели, с учётом их цикличности, в формате int (понедельник 1 ближе к воскресенью 7, чем пятница 5)
func daysBetweenWD(from, to int) int {
	week := make(map[int]int)
	for i := from; i <= 7; {
		if i+1 < 7 {
			week[i+1] = i + 2
		} else {
			week[i+1] = i + 2 - 7
		}
	}
	daysCount := 0
	i := week[from]
	for {
		if daysCount == 7 {
			return daysCount
		}
		if i == to {
			daysCount++
			return daysCount
		} else {
			daysCount++
			i = week[i]
		}
	}
}

// closestWD возвращает ближайший к текущему дню недели, день недели из списка, в формате int, с учётом цикличности недель.
func closestWD(now time.Time, targetDays []int) int {
	if len(targetDays) < 1 {
		return -1
	}
	closestDay := targetDays[0]
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
