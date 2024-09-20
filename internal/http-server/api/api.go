package api

import "todo-list/internal/tasks"

type ErrResponse struct {
	ErrDescription string `json:"error"`
}

func NewErrResponse(description string) ErrResponse {
	return ErrResponse{ErrDescription: description}
}

type PostSucceedResponse struct {
	Id int `json:"id"`
}

func NewPostSucceedResponse(id int) PostSucceedResponse {
	return PostSucceedResponse{Id: id}
}

type TasksResponse struct {
	TaskList []tasks.Task `json:"tasks"`
}
