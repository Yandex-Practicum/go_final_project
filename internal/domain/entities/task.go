// internal/domain/entities/task.go

package entities

// Task представляет задачу в системе
type Task struct {
	ID      int64  `json:"id,omitempty"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}
