package tasks

import (
	"fmt"
	"strconv"
)

type Task struct {
	ID      string `json:"id,omitempty"`
	Date    string `json:"date,omitempty"`
	Title   string `json:"title,omitempty"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

// IDValidation - проверка id на числовое значение
func IDValidation(id string) (Task, error) {

	var Task Task

	_, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("Id не является числом ", err)
		return Task, err
	}

	Task.ID = id

	return Task, nil
}
