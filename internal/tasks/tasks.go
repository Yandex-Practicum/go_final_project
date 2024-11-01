package tasks

import (
	"fmt"
	"strconv"
	"time"
	"todo-server/internal/donetaskrepeat"
	"todo-server/internal/settings"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func (t *Task) CheckID() error {
	if t.ID == "" {
		return fmt.Errorf(`{"error":"Не указан индификатор задачи"}`)
	}
	_, err := strconv.ParseInt(t.ID, 10, 32)
	if err != nil {
		return fmt.Errorf(`{"error":"Указан невозможный индификатор задачи"}`)
	}
	return nil
}

func (t *Task) CheckTitle() error {
	if t.Title == "" {
		return fmt.Errorf(`{"error":"Не указан заголовок задачи"}`)
	}
	return nil
}

func (t *Task) CheckData() (time.Time, error) {
	now := time.Now().Format(settings.Template)
	parsedNow, _ := time.Parse(settings.Template, now)
	if t.Date == "" {
		t.Date = now
	}
	parseDate, err := time.Parse(settings.Template, t.Date)
	if err != nil {
		return parsedNow, fmt.Errorf(`{"error":"Дата указана в неверном формате"}`)
	}
	return parseDate, nil
}

func (t *Task) CheckRepeat(parseDate time.Time) (string, error) {
	if t.Repeat != "" {
		nextDate, err := donetaskrepeat.NextDate(time.Now(), t.Date, t.Repeat)
		if err != nil {
			return "", fmt.Errorf(`{"error":"Неверное правило повторения"}`)
		}
		if parseDate.Before(time.Now()) && t.Date != time.Now().Format(settings.Template) {
			t.Date = nextDate
		}
	} else if parseDate.Before(time.Now()) {
		t.Date = time.Now().Format(settings.Template)
	}
	return t.Date, nil
}
