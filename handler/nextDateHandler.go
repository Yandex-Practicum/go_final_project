package handler

import (
	"go_final_project/nextdate"
	"net/http"
	"time"
)

// NextDateHandler вызывает функцию NextDate и возвращает её результат.
func NextDateHandler(w http.ResponseWriter, req *http.Request) {
	param := req.URL.Query()
	now := param.Get("now")
	day := param.Get("date")
	repeat := param.Get("repeat")

	if now == "" || day == "" || repeat == "" {
		http.Error(w, `{"error":"missing parameters"}`, http.StatusBadRequest)
		return
	}

	timeNow, err := time.Parse("20060102", now)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := nextdate.NextDate(timeNow, day, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result))
}
