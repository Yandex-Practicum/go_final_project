package common

type Task struct {
	ID      string `json:"id,omitempty"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

type Response struct {
	ID    int64  `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}
