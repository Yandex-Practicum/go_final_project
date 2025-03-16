package models

type TaskPutDTO struct {
	Id      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type SignDTO struct {
	Password string `json:"password"`
}
