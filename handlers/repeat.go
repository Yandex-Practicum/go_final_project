package handlers

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// daysInMonth(t) возвращает кол-во дней в месяце даты t
func daysInMonth(t time.Time) int {
	month := t.Month()
	days := time.Date(t.Year(), month+1, 0, 0, 0, 0, 0, t.Location()).Day()
	return days
}

// timeDiff(first, sec) возвращает true, если дата first больше, чем sec
// (проверка на то, чтобы день точно сменился, а не просто одна дата была больше другой)
func timeDiff(first, sec time.Time) bool {
	if first.Year() > sec.Year() {
		return true
	} else if first.Year() == sec.Year() {
		if first.Month() > sec.Month() {
			return true
		} else if first.Month() == sec.Month() {
			if first.Day() > sec.Day() {
				return true
			} else {
				return false
			}
		} else {
			return false
		}
	} else {
		return false
	}
}

// NextDate(now, date, repeat) возвращает следующую от date дату по правилу repeat
// гарантируется, что при этом новая дата будет больше указанной в now
func NextDate(now time.Time, date string, repeat string) (string, error) {
	if len(repeat) == 0 {
		return "", fmt.Errorf("[NextDate]: repeat rule can't be empty")
	}

	nowDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", fmt.Errorf("[NextDate]: wrong date: %v", err)
	}

	repeatError := fmt.Errorf("[NextDate]: wrong repeat format")

	switch repeat[0] {
	case 'd':
		if days := strings.Split(repeat, " "); len(days) != 2 {
			return "", repeatError
		} else {
			d, err := strconv.Atoi(days[1])
			if err != nil {
				return "", repeatError
			}
			if d > 400 {
				return "", repeatError
			}

			nowDate = nowDate.AddDate(0, 0, d)
			for !timeDiff(nowDate, now) {
				nowDate = nowDate.AddDate(0, 0, d)
			}
			return nowDate.Format("20060102"), nil
		}
	case 'y':
		if len(repeat) != 1 {
			return "", repeatError
		}
		nowDate = nowDate.AddDate(1, 0, 0)
		for !timeDiff(nowDate, now) {
			nowDate = nowDate.AddDate(1, 0, 0)
		}
		return nowDate.Format("20060102"), nil
	case 'w':
		if days := strings.Split(repeat, " "); len(days) != 2 {
			return "", repeatError
		} else {
			weekdayToOur := map[string]int{
				"Monday":    1,
				"Tuesday":   2,
				"Wednesday": 3,
				"Thursday":  4,
				"Friday":    5,
				"Saturday":  6,
				"Sunday":    7,
			}

			weekdays := strings.Split(days[1], ",")
			wds := make(map[int]struct{})
			for _, wd := range weekdays {
				toAdd, err := strconv.Atoi(wd)
				if err != nil || toAdd <= 0 || toAdd > 7 {
					return "", repeatError
				}
				wds[toAdd] = struct{}{}
			}

			nowDate = nowDate.AddDate(0, 0, 1)
			for {
				if timeDiff(nowDate, now) {
					if _, ok := wds[weekdayToOur[nowDate.Weekday().String()]]; ok {
						break
					}
				}
				nowDate = nowDate.AddDate(0, 0, 1)
			}

			return nowDate.Format("20060102"), nil
		}
	case 'm':
		if days := strings.Split(repeat, " "); len(days) != 2 && len(days) != 3 {
			return "", repeatError
		} else if len(days) == 2 {
			mds := strings.Split(days[1], ",")
			monthDays := make(map[int]struct{})
			last, prelast := false, false
			for _, md := range mds {
				toAdd, err := strconv.Atoi(md)
				if err != nil || toAdd == 0 || toAdd < -2 || toAdd > 31 {
					return "", repeatError
				}
				if toAdd == -1 {
					last = true
				}
				if toAdd == -2 {
					prelast = true
				}
				monthDays[toAdd] = struct{}{}
			}

			nowDate = nowDate.AddDate(0, 0, 1)
			for {
				if timeDiff(nowDate, now) {

					if _, ok := monthDays[nowDate.Day()]; ok {
						break
					}

					if last && nowDate.Day() == daysInMonth(nowDate) {
						break
					}
					if prelast && nowDate.Day() == daysInMonth(nowDate)-1 {
						break
					}
				}
				nowDate = nowDate.AddDate(0, 0, 1)
			}

			return nowDate.Format("20060102"), nil
		} else if len(days) == 3 {
			mds := strings.Split(days[1], ",")
			monthDays := make(map[int]struct{})

			last, prelast := false, false
			for _, md := range mds {
				toAdd, err := strconv.Atoi(md)
				if err != nil || toAdd == 0 || toAdd < -2 || toAdd > 31 {
					return "", repeatError
				}
				if toAdd == -1 {
					last = true
				}
				if toAdd == -2 {
					prelast = true
				}
				monthDays[toAdd] = struct{}{}
			}

			m := strings.Split(days[2], ",")
			months := make(map[int]struct{})
			for _, month := range m {
				toAdd, err := strconv.Atoi(month)
				if err != nil || toAdd <= 0 || toAdd > 12 {
					return "", repeatError
				}
				months[toAdd] = struct{}{}
			}

			nowDate = nowDate.AddDate(0, 0, 1)
			for {
				if timeDiff(nowDate, now) {

					if _, ok := months[int(nowDate.Month())]; ok {
						if _, ok = monthDays[nowDate.Day()]; ok {
							break
						}
						// случаи, когда нужен последний или предпоследний день(-1 или -2)
						if last && nowDate.Day() == daysInMonth(nowDate) {
							break
						}
						if prelast && nowDate.Day() == daysInMonth(nowDate)-1 {
							break
						}
					}
				}
				nowDate = nowDate.AddDate(0, 0, 1)
			}

			return nowDate.Format("20060102"), nil
		}
	default:
		return "", repeatError
	}

	return "", nil
}
