package models

type Response struct {
	Id    int64  `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

type GetTaskResponseDTO struct {
	Id      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type GetTasksResponseDTO struct {
	Tasks []GetTaskResponseDTO `json:"tasks"`
}

type SignResponseDTO struct {
	Token string `json:"token"`
}
