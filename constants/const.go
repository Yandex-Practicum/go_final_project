package constants

// Task структура задач
type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// TaskIdResponse структура идентификатора созданной записи
type TaskIdResponse struct {
	Id int64 `json:"id"`
}

// ErrorResponse структура ошибки
type ErrorResponse struct {
	Error string `json:"error"`
}

// TaskResponse структура отображения добавленных задач
type TaskResponse struct {
	Tasks []Task `json:"tasks"`
}

const (
	// DateFormat формат даты
	DateFormat = "20060102"
	// Limit лимит для отображения задач
	Limit = 50
)

// WeekMap мапа индексов дней недели
var WeekMap = map[int]int{
	1: 1,
	2: 2,
	3: 3,
	4: 4,
	5: 5,
	6: 6,
	0: 7,
}
