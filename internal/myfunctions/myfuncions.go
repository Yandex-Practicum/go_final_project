package myfunctions

import (
	"fmt"
	"time"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {

	format := "20060102"

	dateParsed, err := time.Parse(format, date)
	if err != nil {
		return "", fmt.Errorf("date parse error: %w\n", err)
	}
	r := dateParsed
	fmt.Printf("dateParsed: %v\n", r)

	timeShift := now.AddDate(0, 0, 1)
	formatTimeShift := timeShift.Format(format)
	return formatTimeShift, nil

}
