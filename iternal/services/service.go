package services

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const (
	TimeFormat = "20060102"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("нет параметра repeat")
	}

	start_date, err := time.Parse(TimeFormat, date)
	if err != nil {
		return "", err
	}

	elements := strings.Fields(repeat)
	if elements[0] == "d" && len(elements) < 2 {
		return "", errors.New("неправильный формат")
	}
	if !strings.Contains("yd", elements[0]) {
		return "", errors.New("неправильный формат")
	}

	for {
		if elements[0] == "y" {
			start_date = start_date.AddDate(1, 0, 0)
		} else if elements[0] == "d" {
			element, _ := strconv.Atoi(elements[1])
			if element > 400 {
				return "", errors.New("неправильный формат")
			}
			start_date = start_date.AddDate(0, 0, element)
		}

		if start_date.After(now) || start_date.Equal(now) {
			return start_date.Format(TimeFormat), nil
		}
	}

}
