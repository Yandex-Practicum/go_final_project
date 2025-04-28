package api

import (
	"fmt"
	"net/http"
)

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dstart := r.FormValue("date")
	repeat := r.FormValue("repeat")

	// вызов NextDate
	nowParsed, err := ParseNowOrDefault(nowStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result, err := NextDate(nowParsed, dstart, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprint(w, result)
}
