package usecases

import "time"

type Task interface {
	GetNextDate(now time.Time, date string, repeat string) (string, error)
}
