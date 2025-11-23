// package dates provides date calculations for scheduler
package dates

import (
	"time"

	"go_final_project/constants"
)

// GetTaskRepetitionDate takes repeat rule and task date and provides new repetition date for task such as an error
func GetTaskRepetitionDate(repeat string, date time.Time) (newDate string, err error) {

	now := time.Now()
	if date.Before(now) {
		if repeat == "" || date.Truncate(24*time.Hour) == date.Truncate(24*time.Hour) {
			newDate = time.Now().Format(constants.DateFormat)
		} else {
			dateStr := date.Format(constants.DateFormat)
			nextDate, errNextDate := NextDate(now, dateStr, repeat)
			err = errNextDate
			newDate = nextDate
		}
	}
	return newDate, err
}
