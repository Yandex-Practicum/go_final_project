package handler

import (
	"encoding/json"
	"go_final_project/internal/models"
	"go_final_project/internal/repository"
	"go_final_project/internal/scheduler"
	"log"
	"net/http"
	"time"
)

func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeErrorJSON(w, "ошибка десериализации JSON", http.StatusBadRequest)
		return
	}

	now := time.Now()
	if err := scheduler.ValidateAndProcessTask(&task, now); err != nil {
		writeErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := repository.AddTaskToDB(&task)
	if err != nil {
		writeErrorJSON(w, "ошибка добавления задачи в базу данных", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err = json.NewEncoder(w).Encode(map[string]interface{}{"id": id})
	if err != nil {
		log.Printf("JSON Encoding Error: %v", err)
	}
}

func writeErrorJSON(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	response := map[string]interface{}{"error": message}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Ошибка encoding JSON error response: %v", err)
	}
}
