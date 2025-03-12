package structurs

import (
	_ "github.com/mattn/go-sqlite3"

	_ "modernc.org/sqlite"
)

type Task struct {
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}
type Tasks struct {
	Id      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}
type Empty struct {
	Out string `json:"out,omitempty"`
}
