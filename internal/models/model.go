package models

type Scheduler struct {
	ID      string `db:"id" json:"id"`           // автоинкрементный идентификатор
	Date    string `db:"date" json:"date"`       // дата в формате YYYYMMDD
	Title   string `db:"title" json:"title"`     // заголовок задачи
	Comment string `db:"comment" json:"comment"` // комментарий к задаче
	Repeat  string `db:"repeat" json:"repeat"`   // правила повторений (максимум 128 символов)
}
