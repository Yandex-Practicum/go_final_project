package utils

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
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

func GetFilterTypeAndValue(r *http.Request) (int, string) {
	search := r.URL.Query().Get("search")
	if len(search) == 0 {
		return FilterTypeNone, ""
	}

	searchDate, err := time.Parse("01.02.2006", search)
	if err == nil {
		return FilterTypeDate, searchDate.Format(ParseDateFormat)
	}
	return FilterTypeSearch, fmt.Sprintf("%%%s%%", search)
}
