package api

import (
	"encoding/json"
	"net/http"

	"todo-list/pkg/db"
)

type TasksResponse struct {
	Tasks []*db.Task `json:"tasks"`
}

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(50)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if tasks == nil {
		tasks = make([]*db.Task, 0)
	}

	response := TasksResponse{
		Tasks: tasks,
	}
	writeJSONResponse(w, response, http.StatusOK)
}

// GetTaskHandler обрабатывает GET запросы для /api/task?id={id}
func GetTaskHandler(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSONError(w, "Не указан идентификатор", http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusNotFound)
		return
	}

	writeJSONResponse(w, task, http.StatusOK)
}

// UpdateTaskHandler обрабатывает PUT запросы для /api/task
func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {

	var req TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, "Ошибка десериализации JSON", http.StatusBadRequest)
		return
	}

	if req.ID == "" {
		writeJSONError(w, "Не указан идентификатор задачи", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		writeJSONError(w, "Не указан заголовок задачи", http.StatusBadRequest)
		return
	}

	if err := processDate(&req); err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Repeat != "" {
		if err := validateRepeat(req.Repeat); err != nil {
			writeJSONError(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	task := &db.Task{
		ID:      req.ID,
		Date:    req.Date,
		Title:   req.Title,
		Comment: req.Comment,
		Repeat:  req.Repeat,
	}

	if err := db.UpdateTask(task); err != nil {
		writeJSONError(w, err.Error(), http.StatusNotFound)
		return
	}

	writeJSONResponse(w, map[string]string{}, http.StatusOK)
}

// DeleteTaskHandler обрабатывает DELETE запросы для /api/task?id={id}
func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSONError(w, "Не указан идентификатор", http.StatusBadRequest)
		return
	}

	if err := db.DeleteTask(id); err != nil {
		writeJSONError(w, err.Error(), http.StatusNotFound)
		return
	}

	writeJSONResponse(w, map[string]string{}, http.StatusOK)
}
