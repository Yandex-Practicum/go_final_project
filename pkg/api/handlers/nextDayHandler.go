package handlers

import (
	"log"
	"net/http"
	"proj/pkg/api/service"
	"time"
)

func NextDayHandler(w http.ResponseWriter, r *http.Request) {
	now := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	time, err := time.Parse(service.LayoutDate, now)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Ошибка парсинга полученного из пути параметра date", http.StatusBadRequest)
		return
	}

	res, err := service.NextDate(time, date, repeat)
	if err != nil {
		http.Error(w, "Ошибка вычисления следующей даты", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(res))

}
