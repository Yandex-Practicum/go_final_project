package models

type Task struct {
	Id      string `json:"id" db:"id"`
	Title   string `json:"title" db:"title"`
	Date    string `json:"date" db:"date"`
	Comment string `json:"comment" db:"comment"`
	Repeat  string `json:"repeat" db:"repeat"`
}
