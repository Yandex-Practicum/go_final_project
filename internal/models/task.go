package models

type Task struct {
	ID      int64  `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type Response struct {
	ID    *int64  `json:"id,omitempty"`
	Error *string `json:"error,omitempty"`
}
