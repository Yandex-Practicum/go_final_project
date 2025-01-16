package configs

const (
	DateFormat = "20060102"
	MaxTasks   = 50
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type TasksResponse struct {
	Tasks []Task `json:"tasks"`
}

type Response struct {
	ID    string `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

var ErrorResponse struct {
	Error string `json:"error,omitempty"`
}
