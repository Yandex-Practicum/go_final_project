package models

// NewTask используется при добавлении новой задачи.
type NewTask struct {
	Date    string `json:"date" db:"date"`
	Title   string `json:"title" db:"title"`
	Comment string `json:"comment" db:"comment"`
	Repeat  string `json:"repeat" db:"repeat"`
}

// Task представляет задачу из базы данных.
// Поле ID теперь имеет тег json:"id,string" – тесты ожидают строковое значение.
type Task struct {
	ID      int64  `json:"id,string" db:"id"`
	Date    string `json:"date" db:"date"`
	Title   string `json:"title" db:"title"`
	Comment string `json:"comment" db:"comment"`
	Repeat  string `json:"repeat" db:"repeat"`
}
