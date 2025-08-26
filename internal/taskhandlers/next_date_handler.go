package taskhandlers

import (
	"go_final_project/internal/domain/services"
	"log"
	"net/http"
	"time"
)

// NextDateHandler handles requests to /api/nextdate
func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now") // Using FormValue() to get GET parameters
	if nowStr == "" {
		nowStr = time.Now().Format("20060102") // If "now" is not provided, use the current date
	}
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	log.Printf("Incoming request to /api/nextdate: now=%s, date=%s, repeat=%s", nowStr, date, repeat)

	if date == "" || repeat == "" {
		log.Println("Error: missing date and repeat parameters")
		http.Error(w, "Date and repeat parameters are required", http.StatusBadRequest)
		return
	}

	nowTime, err := time.Parse("20060102", nowStr)
	if err != nil {
		log.Printf("Error parsing now (%s): %v", nowStr, err)
		http.Error(w, "Invalid date format for now", http.StatusBadRequest)
		return
	}

	nextDate, err := services.NextDate(nowTime, date, repeat)
	if err != nil {
		log.Printf("Error calculating next date: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Successful calculation: nextDate=%s", nextDate)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nextDate))

	log.Printf("Response sent to client: %s", nextDate)
}
