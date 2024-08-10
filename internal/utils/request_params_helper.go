package utils

import (
	"net/http"
	"strconv"
)

func GetIDFromQuery(r *http.Request) (int64, error) {
	stringId := r.URL.Query().Get("id")
	if len(stringId) == 0 {
		return 0, ErrIDIsEmpty
	}
	id, err := strconv.ParseInt(stringId, 10, 64)
	if err != nil {
		return 0, ErrIDIsEmpty
	}
	return id, nil
}
