package entities

type Task struct {
	ID      string  `db:"id" json:"id,omitempty"`
	Date    string  `db:"date" json:"date,omitempty"`
	Title   string  `db:"title" json:"title,omitempty"`
	Comment *string `db:"comment" json:"comment,omitempty"`
	Repeat  string  `db:"repeat" json:"repeat,omitempty"`
}
type Filter struct {
	Limit  int64
	Offset int64
}
type EmptyResponse struct{}

type UserData struct {
	Password string  `json:"password"`
	Token    *string `json:"token"`
}
