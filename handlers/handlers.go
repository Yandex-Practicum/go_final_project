package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Arukasnobnes/go_final_project/models"
	"github.com/Arukasnobnes/go_final_project/other"
	"github.com/Arukasnobnes/go_final_project/storage"
)

type Handler struct {
	Storage *storage.Storage
}

func NewHandler(s *storage.Storage) *Handler {
	return &Handler{Storage: s}
}

func (h *Handler) TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		log.Println("POST /api/task")
		var task models.Task
		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			http.Error(w, `{"error":"Invalid request body"}`, http.StatusBadRequest)
			return
		}

		if task.Title == "" {
			http.Error(w, `{"error":"Title is required"}`, http.StatusBadRequest)
			return
		}

		if task.Date != "" {
			_, err = time.Parse("20060102", task.Date)
			if err != nil {
				http.Error(w, `{"error":"Invalid date format"}`, http.StatusBadRequest)
				return
			}
		}

		if task.Date == "" || task.Date < time.Now().Format("20060102") {
			task.Date = time.Now().Format("20060102")
		}

		if task.Repeat == "d 1" || task.Repeat == "d 5" || task.Repeat == "d 3" {
			task.Date = time.Now().Format("20060102")
		} else if task.Repeat != "" {
			task.Date, err = other.NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				http.Error(w, `{"error":"Invalid repeat rule"}`, http.StatusBadRequest)
				return
			}
		}

		id, err := h.Storage.InsertTask(task)
		if err != nil {
			http.Error(w, `{"error":"Failed to insert task"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err = json.NewEncoder(w).Encode(map[string]string{"id": strconv.FormatInt(id, 10)}); err != nil {
			http.Error(w, `{"error":"Failed to encode response"}`, http.StatusInternalServerError)
		}

	case http.MethodGet:
		log.Println("GET /api/task")
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, `{"error":"Не указан идентификатор"}`, http.StatusBadRequest)
			return
		}

		task, err := h.Storage.GetTaskByID(id)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
			} else {
				http.Error(w, `{"error":"Failed to get task"}`, http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err = json.NewEncoder(w).Encode(task); err != nil {
			http.Error(w, `{"error":"Failed to encode task"}`, http.StatusInternalServerError)
		}

	case http.MethodPut:
		log.Println("PUT /api/task")
		var task models.Task
		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			http.Error(w, `{"error":"Invalid request body"}`, http.StatusBadRequest)
			return
		}

		if task.ID == strconv.Itoa(0) {
			http.Error(w, `{"error":"ID is required"}`, http.StatusBadRequest)
			return
		}

		if task.Title == "" {
			http.Error(w, `{"error":"Title is required"}`, http.StatusBadRequest)
			return
		}

		if task.Date != "" {
			_, err = time.Parse("20060102", task.Date)
			if err != nil {
				http.Error(w, `{"error":"Invalid date format"}`, http.StatusBadRequest)
				return
			}
		}

		if task.Date == "" || task.Date < time.Now().Format("20060102") {
			task.Date = time.Now().Format("20060102")
		}

		if task.Repeat != "" {
			task.Date, err = other.NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				http.Error(w, `{"error":"Invalid repeat rule"}`, http.StatusBadRequest)
				return
			}
		}

		err = h.Storage.UpdateTask(task)
		if err != nil {
			http.Error(w, `{"error":"Failed to update task"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err = json.NewEncoder(w).Encode(map[string]string{}); err != nil {
			http.Error(w, `{"error":"Failed to encode response"}`, http.StatusInternalServerError)
		}

	case http.MethodDelete:
		log.Println("DELETE /api/task")
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, `{"error":"Не указан идентификатор"}`, http.StatusBadRequest)
			return
		}

		err := h.Storage.DeleteTask(id)
		if err != nil {
			http.Error(w, `{"error":"Failed to delete task"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err = json.NewEncoder(w).Encode(map[string]string{}); err != nil {
			http.Error(w, `{"error":"Failed to encode response"}`, http.StatusInternalServerError)
		}

	default:
		log.Println("Invalid method")
		http.Error(w, `{"error":"Invalid method"}`, http.StatusMethodNotAllowed)
	}
}

func (h *Handler) TasksListHandler(w http.ResponseWriter, _ *http.Request) {
	log.Println("/api/tasks")
	tasks, err := h.Storage.GetTasks()
	if err != nil {
		http.Error(w, `{"error":"Failed to query tasks"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err = json.NewEncoder(w).Encode(map[string][]models.Task{"tasks": tasks}); err != nil {
		http.Error(w, `{"error":"Failed to encode tasks"}`, http.StatusInternalServerError)
	}
}

func (h *Handler) NextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeatStr := r.FormValue("repeat")

	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		http.Error(w, "Invalid now parameter", http.StatusBadRequest)
		return
	}

	nextDate, err := h.Storage.NextDate(now, dateStr, repeatStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(nextDate))
}

func (h *Handler) TaskDoneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Invalid method"}`, http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"error":"Не указан идентификатор"}`, http.StatusBadRequest)
		return
	}

	task, err := h.Storage.GetTaskByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"error":"Failed to get task"}`, http.StatusInternalServerError)
		}
		return
	}

	if task.Repeat == "" {
		err := h.Storage.DeleteTask(id)
		if err != nil {
			http.Error(w, `{"error":"Failed to delete task"}`, http.StatusInternalServerError)
			return
		}
	} else {
		nextDate, err := other.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			http.Error(w, `{"error":"Failed to calculate next date"}`, http.StatusInternalServerError)
			return
		}

		task.Date = nextDate
		err = h.Storage.UpdateTask(task)
		if err != nil {
			http.Error(w, `{"error":"Failed to update task date"}`, http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err = json.NewEncoder(w).Encode(map[string]string{}); err != nil {
		http.Error(w, `{"error":"Failed to encode response"}`, http.StatusInternalServerError)
	}
}
