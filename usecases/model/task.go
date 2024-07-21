package model

type TaskReq struct {
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat"`
}

type TaskResp struct {
	Id int64 `json:"id"`
}

func NewTaskResp(id int64) *TaskResp {
	return &TaskResp{Id: id}
}
