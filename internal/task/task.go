package task

import (
	"fmt"
	"time"
)

type Task struct {
	Id      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func Check(t *Task) error {
	if t.Title == "" {
		return fmt.Errorf("Заголовок пустой")
	}

	if t.Date == "" {
		t.Date = time.Now().Format("20060102")
	}

	validDate, err := time.Parse("20060102", t.Date)
	if err != nil {
		return fmt.Errorf("дата должна быть в формате 20060102")
	}

	if t.Repeat != "" && t.Repeat[0] != 'd' && t.Repeat[0] != 'w' && t.Repeat[0] != 'm' && t.Repeat[0] != 'y' {
		return fmt.Errorf("некорректное правило")
	}

	if len(t.Repeat) > 0 {
		if t.Repeat[0] != 'd' && t.Repeat[0] != 'w' && t.Repeat[0] != 'm' && t.Repeat[0] != 'y' {
			return fmt.Errorf("некорректное правило")
		}
		if t.Repeat[0] == 'd' || t.Repeat[0] == 'w' || t.Repeat[0] == 'm' {
			if len(t.Repeat) < 3 {
				return fmt.Errorf("некорректное правило")
			}
		}
	}

	if validDate.Truncate(24 * time.Hour).Before(time.Now().Truncate(24 * time.Hour)) {
		if t.Repeat == "" {
			t.Date = time.Now().Format("20060102")
		}

	}

	if validDate.Truncate(24 * time.Hour).Before(time.Now().Truncate(24 * time.Hour)) {
		if t.Repeat != "" {
			t.Date, err = NextDate(time.Now(), t.Date, t.Repeat)
			if err != nil {
				return fmt.Errorf("нет след даты %v", err)
			}
		}
	}

	return nil
}
