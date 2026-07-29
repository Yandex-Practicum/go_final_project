package service

import "time"

func CheckDate(task *Task) error {
	now := time.Now()
	var next string
	if task.Date == "" {
		task.Date = now.Format(LayoutDate) // если дата не указана берем текущую
	}

	t, err := time.Parse(LayoutDate, task.Date)

	if err != nil {
		return err
	}
	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	nowTester := now.Format(LayoutDate)
	nowTest, err := time.Parse(LayoutDate, nowTester)

	if t.Before(nowTest) {
		if len(task.Repeat) == 0 {
			task.Date = now.Format(LayoutDate)
		} else {
			task.Date = next
		}
	}
	return nil
}
