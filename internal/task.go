package task

import (
	"errors"
	"time"
)

type Task struct {
	ID      string `json:"id,omitempty"`
	Date    string `json:"date,omitempty"`
	Title   string `json:"title,omitempty"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

type List struct {
	Task []*Task `json:"tasks"`
}

const FormatDate = "20060102"

func (t *Task) ValidateForCreate(dateNow string) error {
	if t.Title == "" {
		return errors.New("title is empty")
	}
	if t.Date == "" {
		t.Date = dateNow
		return nil
	}
	_, err := time.Parse(FormatDate, t.Date)
	if err != nil {
		return errors.New("date is invalid")
	}
	if t.Date < dateNow && t.Repeat == "" {
		t.Date = dateNow
	}

	return nil
}
