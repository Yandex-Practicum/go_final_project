package models

type Response struct {
	ID    *int64  `json:"id,omitempty"`
	Error *string `json:"error,omitempty"`
}
