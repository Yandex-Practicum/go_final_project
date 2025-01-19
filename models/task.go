package models

import "time"

// Task представляет задачу
type Task struct {
	ID        string    `json:"id"`
	Date      time.Time `json:"date"`
	Title     string    `json:"title"`
	Comment   string    `json:"comment"`
	Repeat    string    `json:"repeat"`
	Completed bool      `json:"completed"`
}
