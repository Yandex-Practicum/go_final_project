package api

type ErrResponse struct {
	ErrDescription string `json:"error"`
}

type PostSucceedResponse struct {
	Id int `json:"id"`
}

func NewErrResponse(description string) ErrResponse {
	return ErrResponse{ErrDescription: description}
}

func NewPostSucceedResponse(id int) PostSucceedResponse {
	return PostSucceedResponse{Id: id}
}
