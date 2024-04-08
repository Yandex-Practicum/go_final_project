package model

type Sign struct {
	Password string `json:"password"`
}

type AuthToken struct {
	Token string `json:"token"`
}
