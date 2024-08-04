package task

import (
	"fmt"
	"log"
	"time"
)

func Check(t *Task) error {
	log.Println("checking task")

	if t.Title == "" {
		return fmt.Errorf("title is empty")
	}

	if t.Date == "" {
		t.Date = time.Now().Format(formatOfDate)
	}

	validDate, err := time.Parse(formatOfDate, t.Date)
	if err != nil {
		return fmt.Errorf("date is in a format other than 20060102")
	}

	if t.Repeat != "" && t.Repeat[0] != 'd' && t.Repeat[0] != 'w' && t.Repeat[0] != 'm' && t.Repeat[0] != 'y' {
		return fmt.Errorf("incorrect repetition rule")
	}

	if len(t.Repeat) > 0 {
		if t.Repeat[0] != 'd' && t.Repeat[0] != 'w' && t.Repeat[0] != 'm' && t.Repeat[0] != 'y' {
			return fmt.Errorf("incorrect repetition rule")
		}
		if t.Repeat[0] == 'd' || t.Repeat[0] == 'w' || t.Repeat[0] == 'm' {
			if len(t.Repeat) < 3 {
				return fmt.Errorf("incorrect repetition rule")
			}
		}
	}

	if validDate.Truncate(24 * time.Hour).Before(time.Now().Truncate(24 * time.Hour)) {
		if t.Repeat == "" {
			t.Date = time.Now().Format(formatOfDate)
		}
	}

	if validDate.Truncate(24 * time.Hour).Before(time.Now().Truncate(24 * time.Hour)) {
		if t.Repeat != "" {
			t.Date, err = NextDate(time.Now(), t.Date, t.Repeat)
			if err != nil {
				return fmt.Errorf("can't get nextDate: %v", err)
			}
		}
	}

	return nil
}
