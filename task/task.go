package task

import (
	"errors"
	"final/nextdate"
	"time"
)

const ParseDate = "20060102"

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// CheckTitle проверяет наличие заголовка задачи.
func (t Task) CheckTitle() error {
	if t.Title == "" {
		return errors.New("пустой заголовок")
	}
	return nil
}

// CheckDate проверяет корректность даты задачи и устанавливает ее в настоящее время, если она неправильная.
func (t Task) CheckDate() (Task, error) {
	now := time.Now()
	if t.Date == "" {
		t.Date = now.Format(ParseDate)
		return t, nil
	}

	date, err := time.Parse(ParseDate, t.Date)
	if err != nil {
		return t, errors.New("неправильный формат даты")
	}
	if date.Before(now) {
		if t.Repeat == "" {
			t.Date = now.Format(ParseDate)
			return t, nil
		}
		nextDate, err := nextdate.CalcNextDate(now.Format(ParseDate), t.Date, t.Repeat)
		if err != nil {
			return t, errors.New("ошибка вычисления следующей даты")
		}
		t.Date = nextDate
	}
	return t, nil
}

// CountDate вычисляет следующую дату, если задано правило повторения.
func (t Task) CountDate() error {
	if t.Repeat != "" {
		nextDate, err := nextdate.CalcNextDate(time.Now().Format(ParseDate), t.Date, t.Repeat)
		if err != nil {
			return errors.New("ошибка вычисления следующей даты")
		}
		t.Date = nextDate
	}
	return nil
}

// CheckID проверяет наличие идентификатора задачи.
func (t Task) CheckID() error {
	if t.ID == "" {
		return errors.New("не указан идентификатор задачи")
	}
	return nil
}

// CheckRepeat проверяет корректность правила повторения.
func (t Task) CheckRepeat() error {
	if t.Repeat != "" {
		if _, err := nextdate.ParseRepeatRules(t.Repeat); err != nil {
			return errors.New("правило повторения указано в неправильном формате")
		}
	}
	return nil
}
