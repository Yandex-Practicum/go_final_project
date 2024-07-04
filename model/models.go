package model

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type ResponseError struct {
	Error string `json:"error"`
}

type Tasks struct {
	Tasks []Task `json:"tasks"`
}
