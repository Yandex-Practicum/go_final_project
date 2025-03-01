package myfunctions

import (
	"time"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {

	nextDate := now.AddDate(0, 0, 1)
	formattedNextDate := nextDate.Format("20060102")
	return formattedNextDate, nil
	
}
