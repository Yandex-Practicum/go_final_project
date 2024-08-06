// package handlers provides API handlers
package handlers

import (
	"errors"
	"fmt"
	"go_final_project/tasks"
	"log"
	"net/http"
	"time"
)

// GetNextDate takes request and calculate new repetition date for task
func GetNextDate(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request GetNextDate")

	nowStr := r.FormValue("now")
	if nowStr == "" {
		err := errors.New("missing 'now' parameter")
		log.Printf("Missing 'now' parameter: %v", err)
		http.Error(w, "Missing 'now' parameter", http.StatusBadRequest)
		return
	}

	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		log.Printf("Invalid 'now' date format: %v", err)
		http.Error(w, fmt.Sprintf("Invalid 'now' date format: %v", err), http.StatusBadRequest)
		return
	}

	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	result, err := tasks.NextDate(now, date, repeat)
	if err != nil {
		log.Printf("Error calculating next date: %v", err)
		http.Error(w, fmt.Sprintf("Error calculating next date: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write([]byte(result)); err != nil {
		log.Printf("Error writing response: %v\n", err)
	}
}
