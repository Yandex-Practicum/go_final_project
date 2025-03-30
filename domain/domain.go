package domain

type Task struct {
	ID      string `json:"id,omitempty" `
	Date    string `json:"date" `
	Title   string `json:"title" `
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat" `
}
