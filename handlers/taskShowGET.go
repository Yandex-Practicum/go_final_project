package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/LEbauchoir/go_final_project/models"
)

func TasksShowGET(w http.ResponseWriter, r *http.Request) {
	var tasks []models.Task
	var err error

	if tasks, err = dbHelper.TasksShow(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Ошибка при выводе Tasks из БД: %v", err)
		return
	}

	if tasks == nil {
		tasks = []models.Task{}
	}

	tasksData, err := json.Marshal(models.Tasks{Tasks: tasks})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error marshaling tasksData: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(tasksData)
	if err != nil {
		log.Printf("Ошибка при ответе: %v", err)
	}
}
