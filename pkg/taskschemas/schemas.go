package taskschemas

// TaskSchema описывает структуру задачи в API
type TaskSchema struct {
	ID      int64  `json:"id,omitempty"` // ID задачи (может отсутствовать для новых задач)
	Date    string `json:"date"`         // Дата задачи в формате YYYYMMDD
	Title   string `json:"title"`        // Заголовок задачи
	Comment string `json:"comment"`      // Комментарий к задаче
	Repeat  string `json:"repeat"`       // Правило повторения (может быть пустым)
}
