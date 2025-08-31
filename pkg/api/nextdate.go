package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const DateFmt = "20060102"

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	nowStr := strings.TrimSpace(r.FormValue("now"))
	dateStr := strings.TrimSpace(r.FormValue("date"))
	repeat := strings.TrimSpace(r.FormValue("repeat"))

	var now time.Time
	if nowStr == "" {
		now = today(time.Now())
	} else {
		t, err := parseDate(nowStr)
		if err != nil {
			_, _ = w.Write([]byte("bad now\n"))
			return
		}
		now = t
	}

	if dateStr == "" || repeat == "" {
		_, _ = w.Write([]byte("missing date or repeat\n"))
		return
	}

	next, err := NextDate(now, dateStr, repeat)
	if err != nil {
		_, _ = w.Write([]byte(err.Error() + "\n"))
		return
	}
	_, _ = w.Write([]byte(next + "\n"))
}

func NextDate(now time.Time, startStr string, repeat string) (string, error) {
	start, err := parseDate(startStr)
	if err != nil {
		return "", errors.New("bad date")
	}
	repeat = strings.ToLower(strings.TrimSpace(repeat))
	if repeat == "" {
		return "", errors.New("empty repeat")
	}

	now = today(now)

	if strings.HasPrefix(repeat, "d") {
		n, ok, e := parseDaily(repeat)
		if e != nil {
			return "", e
		}
		if !ok {
			return "", errors.New("bad repeat")
		}
		t := start.AddDate(0, 0, n)
		for !afterDay(t, now) {
			t = t.AddDate(0, 0, n)
		}
		return t.Format(DateFmt), nil
	}

	if repeat == "y" {
		t := addYear(start)
		for !afterDay(t, now) {
			t = addYear(t)
		}
		return t.Format(DateFmt), nil
	}

	if strings.HasPrefix(repeat, "w") {
		week, ok, e := parseWeekly(repeat)
		if e != nil {
			return "", e
		}
		if !ok {
			return "", errors.New("bad repeat")
		}
		t := start.AddDate(0, 0, 1)
		steps := 0
		for {
			if inWeek(t, week) && afterDay(t, now) {
				return t.Format(DateFmt), nil
			}
			t = t.AddDate(0, 0, 1)
			steps++
			if steps > 366*10 {
				return "", errors.New("no next date found")
			}
		}
	}

	if strings.HasPrefix(repeat, "m") {
		days, months, ok, e := parseMonthly(repeat)
		if e != nil {
			return "", e
		}
		if !ok {
			return "", errors.New("bad repeat")
		}
		t := start.AddDate(0, 0, 1)
		steps := 0
		for {
			if monthOk(t, days, months) && afterDay(t, now) {
				return t.Format(DateFmt), nil
			}
			t = t.AddDate(0, 0, 1)
			steps++
			if steps > 366*15 {
				return "", errors.New("no next date found")
			}
		}
	}

	return "", errors.New("unsupported repeat")
}

func parseDate(s string) (time.Time, error) {
	t, err := time.Parse(DateFmt, s)
	if err != nil {
		return time.Time{}, err
	}
	return today(t), nil
}

func today(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

func afterDay(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	if ay != by {
		return ay > by
	}
	if am != bm {
		return am > bm
	}
	return ad > bd
}

func parseDaily(s string) (int, bool, error) {
	parts := strings.Fields(s)
	if len(parts) != 2 || parts[0] != "d" {
		return 0, false, errors.New("bad repeat")
	}
	n, err := strconv.Atoi(parts[1])
	if err != nil || n < 1 || n > 400 {
		return 0, false, errors.New("bad daily interval")
	}
	return n, true, nil
}

func addYear(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y+1, m, d, 0, 0, 0, 0, t.Location())
}

func parseWeekly(s string) ([8]bool, bool, error) {
	var week [8]bool
	parts := strings.Fields(s)
	if len(parts) != 2 || parts[0] != "w" {
		return week, false, errors.New("bad repeat")
	}
	items := strings.Split(parts[1], ",")
	if len(items) == 0 {
		return week, false, errors.New("bad weekly")
	}
	for _, it := range items {
		it = strings.TrimSpace(it)
		n, err := strconv.Atoi(it)
		if err != nil || n < 1 || n > 7 {
			return week, false, errors.New("bad weekly")
		}
		week[n] = true
	}
	return week, true, nil
}

func inWeek(t time.Time, week [8]bool) bool {
	w := int(t.Weekday())
	idx := (w+6)%7 + 1
	return week[idx]
}

func parseMonthly(s string) ([]int, []int, bool, error) {
	parts := strings.Fields(s)
	if len(parts) < 2 || len(parts) > 3 || parts[0] != "m" {
		return nil, nil, false, errors.New("bad repeat")
	}

	var days []int
	for _, it := range strings.Split(parts[1], ",") {
		it = strings.TrimSpace(it)
		if it == "" {
			continue
		}
		n, err := strconv.Atoi(it)
		if err != nil {
			return nil, nil, false, errors.New("bad monthly days")
		}
		if !(n >= 1 && n <= 31) && n != -1 && n != -2 {
			return nil, nil, false, errors.New("bad monthly days")
		}
		days = append(days, n)
	}
	if len(days) == 0 {
		return nil, nil, false, errors.New("bad monthly days")
	}

	var months []int
	if len(parts) == 3 {
		for _, it := range strings.Split(parts[2], ",") {
			it = strings.TrimSpace(it)
			if it == "" {
				continue
			}
			n, err := strconv.Atoi(it)
			if err != nil || n < 1 || n > 12 {
				return nil, nil, false, errors.New("bad monthly months")
			}
			months = append(months, n)
		}
		if len(months) == 0 {
			return nil, nil, false, errors.New("bad monthly months")
		}
	}
	return days, months, true, nil
}

func monthOk(t time.Time, days []int, months []int) bool {
	if len(months) > 0 {
		cur := int(t.Month())
		ok := false
		for _, m := range months {
			if m == cur {
				ok = true
				break
			}
		}
		if !ok {
			return false
		}
	}
	ld := lastDay(t)
	d := t.Day()
	for _, want := range days {
		if want > 0 && want == d && want <= ld {
			return true
		}
		if want == -1 && d == ld {
			return true
		}
		if want == -2 && d == ld-1 {
			return true
		}
	}
	return false
}

func lastDay(t time.Time) int {
	y, m, _ := t.Date()
	firstNext := time.Date(y, m+1, 1, 0, 0, 0, 0, t.Location())
	return firstNext.AddDate(0, 0, -1).Day()
}
