package main

import (
	"fmt"
	"net/http"
	"time"
)

func getNextDate(w http.ResponseWriter, r *http.Request) {

	format := "20060102"

	q := r.URL.Query()
	now := q.Get("now")
	date := q.Get("date")
	repeat := q.Get("repeat")

	nowDate, err := time.Parse(format, now)
	if err != nil {
		fmt.Println(err)
		return
	}

	nextDate, err := NextDate(nowDate, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := []byte(nextDate)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
