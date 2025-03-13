package dateutil

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const _dateLayout = "20060102"

type RepeatRule string

const (
	daily RepeatRule = "d"
	// weekly  RepeatRule = "w"
	// monthly RepeatRule = "m"
	yearly RepeatRule = "y"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("expected repeat, got an empty string")
	}

	_, err := strconv.Atoi(date)
	if err != nil {
		return "", fmt.Errorf("wrong date type: %w", err)
	}

	dt, err := time.Parse(_dateLayout, date)
	if err != nil {
		return "", fmt.Errorf("unable to parse date: %w", err)
	}

	nr, err := nextRepeat(repeat)
	if err != nil {
		return "", err
	}

	return dt.AddDate(nr.Year, nr.Month, nr.Day).Format(_dateLayout), nil
}

type repeat struct {
	Year, Month, Day int
}

func nextRepeat(rule string) (*repeat, error) {
	var next = &repeat{}

	s := strings.Split(rule, " ")

	rr := RepeatRule(s[0])

	switch rr {
	case daily:
		var n int
		n, err := strconv.Atoi(s[1])
		if err != nil {
			return nil, err
		}
		next.Day = n

	case yearly:
		next.Year = 1

	default:
		return nil, fmt.Errorf("unkown repeat identifier: %q", rule)
	}

	return next, nil
}
