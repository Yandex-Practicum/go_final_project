package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/FunnyFoXD/go_final_project/databases"
	"github.com/FunnyFoXD/go_final_project/helpers"
	"github.com/FunnyFoXD/go_final_project/models"
)

type insertID struct {
	ID int `json:"id"`
}

func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	var id insertID
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, `{"error":"can't read body"}`, http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, `{"error":"can't unmarshal body"}`, http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		http.Error(w, `{"error":"title is empty"}`, http.StatusBadRequest)
		return
	}

	log.Println(task)

	if task.Date == "" || task.Date == "today" || task.Date == "Today" {
		task.Date = time.Now().Format("20060102")
	}

	taskParse, err := time.Parse("20060102", task.Date)
	if err != nil {
		http.Error(w, `{"error":"invalid date format"}`, http.StatusBadRequest)
		return
	}

	log.Println(taskParse.Before(time.Now()))
	if taskParse.Before(time.Now()) {
		if task.Repeat == "" {
			task.Date = time.Now().Format("20060102")
		} else {
			nextDate, err := helpers.NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusBadRequest)
				return
			}

			task.Date = nextDate
		}
	}

	id.ID, err = databases.InsertTask(task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
