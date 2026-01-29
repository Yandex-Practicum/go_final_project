package common

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go_final_project/internal/constants"
)

func GetIdFromQuery(r *http.Request) (int64, error) {
	stringId := r.URL.Query().Get("id")
	if len(stringId) == 0 {
		return 0, constants.ErrIdIsEmpty
	}
	id, err := strconv.ParseInt(stringId, 10, 64)
	if err != nil {
		return 0, constants.ErrIdIsEmpty
	}
	return id, nil
}

func GetFilterTypeAndValue(r *http.Request) (int, string) {
	search := r.URL.Query().Get("search")
	if len(search) == 0 {
		return constants.FilterTypeNone, ""
	}

	searchDate, err := time.Parse("02.01.2006", search)
	if err == nil {
		return constants.FilterTypeDate, searchDate.Format(constants.ParseDateFormat)
	}
	return constants.FilterTypeSearch, fmt.Sprintf("%%%s%%", search)
}
