package model

type Task struct {
	Id      string `json:"id,omitempty"` // идентификатор задачи
	Date    string `json:"date"`         // дата задачи в формате "20060102"
	Title   string `json:"title"`        // наименование задачи
	Comment string `json:"comment"`      // комментарий
	Repeat  string `json:"repeat"`       // правило повторения задачи
}

type TaskID struct {
	Id string `json:"id"`
}

type TaskError struct {
	Error string `json:"error,omitempty"`
}

type TasksType struct {
	Tasks []Task `json:"tasks"`
}
