package handler

import (
	"log"
	"net/http"
	"time"

	"go_final_project/task"
)

func GetNextDate(w http.ResponseWriter, r *http.Request) {
	log.Println("Received reqest GetNextDate")

	r.ParseForm()

	now, err := time.Parse("20060102", r.FormValue("now"))
	if err != nil {
		log.Printf("Incorrect now date: %v", err)
		http.Error(w, "Incorrect now date", http.StatusBadRequest)
		return
	}

	date := r.FormValue("date")

	repeat := r.FormValue("repeat")

	result, err := task.NextDate(now, date, repeat)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "string")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write([]byte(result))
	if err != nil {
		log.Println("Error write in func GetNextDate:", err)
	}
}
