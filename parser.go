package parser

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

type DRepeat struct {
	num int
}


// ParseDRepeat заполняет структуру DRepeat
func ParseDRepeat(rule []string) (*DRepeat, error) {
	next, err := strconv.Atoi(rule[1])
	if err != nil {
		return nil, fmt.Errorf("Error in checking days in repeat rule, got '%s'", rule[1])
	}
	if next > 0 && next <= 400 {
		return &DRepeat{num: next}, nil
	}
	return nil, fmt.Errorf("expected number of days less than 400, got '%s'", rule[1])
}

func (dr *DRepeat) GetNextDate(now time.Time, date time.Time) (time.Time, error) {
	result := date
	for {
		result = result.AddDate(0, 0, dr.num)
		if result.After(now) {
			return result, nil
		}
	}
}

// --------------------------------------------------------
type YRepeat struct {
}

// ParseYRepeat заполняет структуру YRepeat
func ParseYRepeat(rule []string) (*YRepeat, error) {
	return &YRepeat{}, nil
}

func (yr *YRepeat) GetNextDate(now time.Time, date time.Time) (time.Time, error) {
	i := 1

	for {
		result := date.AddDate(i, 0, 0)
		if result.After(now) {
			return result, nil
		}
		i++
	}
}

// --------------------------------------------------------
type WRepeat struct {
	nums []int
}

// ParseWRepeat заполняет структуру WRepeat
func ParseWRepeat(rule []string) (*WRepeat, error) {
	// var week: days of week in rule repeat
	if len(rule) == 1 {
		return nil, fmt.Errorf("Error in w rule.")
	}

	week := []int{}
	x := strings.Split(rule[1], ",")
	for i := 0; i < len(x); i++ {
		num, err := strconv.Atoi(x[i])
		if err != nil || num > 7 || num < 1 {
			return nil, fmt.Errorf("Can not parse days for repeat value.")
		}
		week = append(week, num)
	}
	return &WRepeat{nums: week}, nil
}

func (wr *WRepeat) GetNextDate(now time.Time, date time.Time) (time.Time, error) {
	startdate := startDateForMWrule(now, date)

	todayWeekday := startdate.Weekday()

	sort.Ints(wr.nums) // сортируем, чтобы сразу взять тот день, что больше номером, чем сегодняшний

	numDay := int(todayWeekday)
	if numDay == 7 {
		numDay = 0
	}

	for _, n := range wr.nums {
		if n > numDay {
			result := startdate.AddDate(0, 0, n-numDay)
			return result, nil
		}
	}

	increment := 7 - int(startdate.Weekday())

	result := startdate.AddDate(0, 0, increment+wr.nums[0])

	return result, nil
}

// --------------------------------------------------------
type MRepeat struct {
	mDays   []int
	mMonths []int
}

// hasMonths определяет, есть ли в правиле месяцы
func (mr *MRepeat) hasMonths() bool {
	return len(mr.mMonths) > 0
}

// ParseMRepeat first checks the the second part of repeat rule string for "m" rule.
// -1 and -2 are converted to the last day of the month and day before last.
// Then checks the the third part of repeat rule string for "m" rule.
func ParseMRepeat(rule []string, now time.Time) (*MRepeat, error) {
	// Сначала всегда рассматриваем сегодняшний месяц.
	// Смотрим со всеми днями, а уж если не подходят предложенные правилом дни (все < сегодня),
	// то месяц берем следующий месяц.
	// Если в правиле 31ое число, а в рассматриваемом месяце 30 дней,
	// то проверяется, чтобы в следующем месяце был 31 день и рассматриваем уже следующий месяц
	// var lenOfMRule int

	if len(rule) == 1 || len(rule) > 3 {
		return nil, fmt.Errorf("Error in m rule.")
	}

	// определяем, есть ли месяцы в правиле
	hasMonths := false

	if len(rule) == 3 {
		hasMonths = true
	}

	days := []int{}

	daysInRule := strings.Split(rule[1], ",") // daysInRule - это дни в правиле m

	for _, day := range daysInRule {
		num, err := strconv.Atoi(day)
		if err != nil {
			return nil, fmt.Errorf("Error in checking days in repeat rule 'm', got '%s'", day)
		}
		if num >= 1 && num <= 31 {
			days = append(days, num)
		} else if num == -1 {
			// time.Date принимает значения вне их обычных диапазонов, то есть
			// значения нормализуются во время преобразования
			// Чтобы рассчитать количество дней текущего месяца (t), смотрим на день следующего месяца
			t := Date(now.Year(), int(now.Month()+1), 0)
			days = append(days, int(t.Day()))
		} else if num == -2 {
			// time.Date принимает значения вне их обычных диапазонов, то есть
			// значения нормализуются во время преобразования
			// Чтобы рассчитать количество дней текущего месяца (t), смотрим на день следующего месяца
			t := Date(now.Year(), int(now.Month()+1), 0)
			days = append(days, int(t.Day())-1)
		} else {
			return nil, fmt.Errorf("Error in checking days in repeat rule 'm', got '%s'", day)
		}

	}

	// checks the the third part of repeat rule string for "m" rule.
	months := []int{}

	if hasMonths {
		monthsInRule := strings.Split(rule[2], ",") // monthsInRule - это месяцы в правиле m

		for _, month := range monthsInRule {
			num, err := strconv.Atoi(month)
			if err != nil || num < 1 || num > 12 {
				return nil, fmt.Errorf("Error in checking days in repeat rule 'm', got '%s'", month)
			}
			months = append(months, num)
		}
	}
	return &MRepeat{mDays: days, mMonths: months}, nil
}

func (mr *MRepeat) GetNextDate(now time.Time, date time.Time) (time.Time, error) {
	startdate := startDateForMWrule(now, date)

	sort.Ints(mr.mDays)

	// ниже проверяем, что день startdate не является больше, чем последнее число из mDays
	// если же больше, то startmonth надо сделать следующим месяцем
	var nextDay time.Time

	if !mr.hasMonths() {
		for _, day := range mr.mDays {
			if day > int(startdate.Day()) {
				nextDay = startdate.AddDate(0, 0, day-int(startdate.Day()))
				if nextDay.Day() != day {
					nextDay = Date(startdate.Year(), int(startdate.Month())+1, day)
				}
				return nextDay, nil
			}
		}

		if nextDay == Date(0001, 1, 1) { // 0001-01-01 00:00:00 +0000 UTC нулевой вариант времени
			startdate = Date(int(startdate.Year()), int(startdate.Month())+1, 1)
			for _, day := range mr.mDays {
				if day >= int(startdate.Day()) {
					nextDay = startdate.AddDate(0, 0, day-int(startdate.Day()))
					return nextDay, nil
				}
			}
		}
	}

	if mr.hasMonths() {
		sort.Ints(mr.mMonths)

		nextDay = ruleMwithMonth(startdate, mr.mDays, mr.mMonths)
		return nextDay, nil
	}

	return time.Time{}, fmt.Errorf("Error in checking days and months in 'm' repeat rule")
}

// ----------------------------------------------------------------

type Repeat interface {
	GetNextDate(now time.Time, date time.Time) (time.Time, error)
}

// ----------------------------------------------------------------

// Date returns time type from the int types of year, month and day.
func Date(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

// startDateForMWrule selects startdate from now and date: selects a later date
func startDateForMWrule(now time.Time, date time.Time) time.Time {
	if date.After(now) {
		return date
	}
	return now
}

// ruleMwithMonth gets sorted mDays and mMonths and returns nextDay
func ruleMwithMonth(startdate time.Time, mDays []int, mMonths []int) time.Time {
	var nextDay time.Time

	for _, month := range mMonths {
		if month == int(startdate.Month()) {
			startdate = Date(startdate.Year(), month, 1)

			// dayInMonth is number of days in the current month.
			t := Date(startdate.Year(), int(startdate.Month())+1, 0) // день до следующего месяца
			dayInMonth := t.Day()

			for _, day := range mDays {
				if day > int(startdate.Day()) && day <= dayInMonth {
					gotDay := Date(startdate.Year(), int(startdate.Month()), day)
					nextDay = gotDay
					return nextDay
				} else if day > int(startdate.Day()) && day > dayInMonth {
					startdate = Date(startdate.Year(), int(startdate.Month())+1, 1)
				}
			}
		} else if month > int(startdate.Month()) { // else сделан для того,
			// чтобы 1 число следующего месяца тоже учитывалось в поиске
			startdate = Date(startdate.Year(), month, 1)

			// dayInMonth is number of days in the current month.
			t := Date(startdate.Year(), int(startdate.Month())+1, 0) // день до следующего месяца
			dayInMonth := t.Day()

			for _, day := range mDays {
				if day >= int(startdate.Day()) && day <= dayInMonth {
					gotDay := Date(startdate.Year(), int(startdate.Month()), day)
					nextDay = gotDay
					return nextDay
				} else if day > int(startdate.Day()) && day > dayInMonth {
					startdate = Date(startdate.Year(), int(startdate.Month())+1, 1)
				}
			}
		}
	}
	return nextDay
}