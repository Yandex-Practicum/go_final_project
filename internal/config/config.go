package config

// TimeFormat задает формат даты
const TimeFormat = "20060102"

// Task представляет задачу
type Task struct {
	Id      string `json:"id,omitempty"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat"`
}
