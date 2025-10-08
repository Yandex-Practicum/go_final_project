package model

import (
	"encoding/json"
	"strconv"
)

type Int64String int64

// Метод для получения значения в формате int64
func (i *Int64String) Int64() int64 {
	return int64(*i)
}

func (i Int64String) MarshalJSON() ([]byte, error) {
	return json.Marshal(strconv.FormatInt(int64(i), 10))
}

type Task struct {
	ID      Int64String `db:"id" json:"id,omitempty"`
	Date    string      `db:"date" json:"date"`
	Title   string      `db:"title" json:"title"`
	Comment string      `db:"comment" json:"comment,omitempty"`
	Repeat  string      `db:"repeat" json:"repeat,omitempty"`
}

type TasksResponse struct {
	Tasks []Task `json:"tasks"`
}
