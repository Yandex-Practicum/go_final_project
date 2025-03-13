package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"go_final-project/internal/auth"
	"go_final-project/internal/db"
	"go_final-project/internal/logic"
	"go_final-project/internal/task"
	"net/http"
	"os"
	"strconv"
	"time"
)

// GetTasksHandler handles retrieving a list of tasks, supporting search and filtering.
func GetTasksHandler(dbase *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			http.Error(w, `{"error":"Only GET method is supported"}`, http.StatusMethodNotAllowed)
			return
		}
		search := req.URL.Query().Get("search")
		var dateSearch string

		if parsedDate, err := time.Parse("02.01.2006", search); err == nil {
			dateSearch = parsedDate.Format("20060102")
		}

		tasks, err := db.GetTasks(dbase, search, dateSearch)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}

		if tasks == nil {
			tasks = []task.Task{}
		}

		tasksList := make([]map[string]interface{}, len(tasks))
		for i, t := range tasks {
			tasksList[i] = map[string]interface{}{
				"id":      strconv.FormatInt(t.ID, 10),
				"date":    t.Date,
				"title":   t.Title,
				"comment": t.Comment,
				"repeat":  t.Repeat,
			}
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(map[string]interface{}{"tasks": tasksList})
	}
}

// TaskHandler handles task operations: GET (retrieve), POST (create), PUT (update), DELETE (remove).
func TaskHandler(dbase *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			getTask(dbase, w, req)
		case http.MethodPost:
			addTask(dbase, w, req)
		case http.MethodPut:
			updateTask(dbase, w, req)
		case http.MethodDelete:
			deleteTask(dbase, w, req)
		default:
			sendJSONError(w, "Only GET, POST, PUT, DELETE methods are supported.", http.StatusMethodNotAllowed)
		}
	}
}

// getTask retrieves a task from the database by its ID.
func getTask(dbase *sqlx.DB, w http.ResponseWriter, req *http.Request) {
	idStr := req.URL.Query().Get("id")
	if idStr == "" {
		sendJSONError(w, "id is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		sendJSONError(w, "id error format", http.StatusBadRequest)
		return
	}

	task, err := db.GetTaskByID(dbase, id)
	if err != nil {
		sendJSONError(w, "issue not found", http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"id":      strconv.FormatInt(task.ID, 10),
		"date":    task.Date,
		"title":   task.Title,
		"comment": task.Comment,
		"repeat":  task.Repeat,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(response)
}

// addTask inserts a new task into the database, ensuring valid input.
func addTask(dbase *sqlx.DB, w http.ResponseWriter, req *http.Request) {
	var newTask task.Task

	err := json.NewDecoder(req.Body).Decode(&newTask)
	if err != nil {
		sendJSONError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if newTask.Title == "" {
		sendJSONError(w, "Title is required", http.StatusBadRequest)
		return
	}
	now := time.Now().Truncate(24 * time.Hour)
	today := now.Format("20060102")

	if newTask.Date == "" {
		newTask.Date = today
	} else {
		if _, err := time.Parse("20060102", newTask.Date); err != nil {
			sendJSONError(w, "date error format", http.StatusBadRequest)
			return
		}
	}
	taskDate, _ := time.Parse("20060102", newTask.Date)

	if newTask.Repeat != "" {
		if taskDate.Before(now) {
			nextDate, err := logic.NextDate(now, newTask.Date, newTask.Repeat)
			if err != nil || nextDate == "" {
				sendJSONError(w, "Invalid repeat format or no valid next date found", http.StatusBadRequest)
				return
			}
			newTask.Date = nextDate
		}
	} else if taskDate.Before(now) {
		newTask.Date = today
	}

	id, err := db.AddTask(dbase, &newTask)
	if err != nil {
		sendJSONError(w, "Failed to save task", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(map[string]string{"id": strconv.FormatInt(id, 10)})
}

// updateTask handles updating an existing task in the database.
// It reads the request body, validates input data, and updates the task if it exists.
// If the task is not found or the input is invalid, it returns an appropriate error response.
func updateTask(dbase *sqlx.DB, w http.ResponseWriter, req *http.Request) {
	var rawData map[string]interface{}
	err := json.NewDecoder(req.Body).Decode(&rawData)
	if err != nil {
		sendJSONError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	idValue, ok := rawData["id"]
	if !ok {
		sendJSONError(w, "id is required", http.StatusBadRequest)
		return
	}

	var id int64
	switch v := idValue.(type) {
	case string:
		id, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			sendJSONError(w, "id error format", http.StatusBadRequest)
			return
		}
	case float64:
		id = int64(v)
	default:
		sendJSONError(w, "id error format", http.StatusBadRequest)
		return
	}

	rawData["id"] = id

	updatedJSON, err := json.Marshal(rawData)
	if err != nil {
		sendJSONError(w, "Server error", http.StatusInternalServerError)
		return
	}

	var updatedTask task.Task
	err = json.Unmarshal(updatedJSON, &updatedTask)
	if err != nil {
		sendJSONError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	existingTask, err := db.GetTaskByID(dbase, id)
	if err != nil {
		sendJSONError(w, "Task not found", http.StatusNotFound)
		return
	}

	if updatedTask.Title == "" {
		sendJSONError(w, "Title is required", http.StatusBadRequest)
		return
	}

	if updatedTask.Date == "" {
		updatedTask.Date = existingTask.Date
	}
	if updatedTask.Comment == "" {
		updatedTask.Comment = existingTask.Comment
	}
	if updatedTask.Repeat == "" {
		updatedTask.Repeat = existingTask.Repeat
	}

	if _, err := time.Parse("20060102", updatedTask.Date); err != nil {
		sendJSONError(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	if updatedTask.Repeat != "" {
		_, err := logic.NextDate(time.Now(), updatedTask.Date, updatedTask.Repeat)
		if err != nil {
			sendJSONError(w, "Invalid repeat format", http.StatusBadRequest)
			return
		}
	}

	err = db.UpdateTask(dbase, &updatedTask)
	if err != nil {
		sendJSONError(w, "issue not found", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(map[string]string{})
}

// MarkTaskDoneHandler marks a task as completed based on its ID.
// - If the task is non-recurring, it is deleted from the database.
// - If the task is recurring, its next execution date is calculated and updated.
// Returns an appropriate JSON response indicating success or failure.
func MarkTaskDoneHandler(dbase *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			sendJSONError(w, "Invalid method", http.StatusMethodNotAllowed)
			return
		}

		idStr := req.URL.Query().Get("id")
		if idStr == "" {
			sendJSONError(w, "id is required", http.StatusBadRequest)
			return
		}

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || id <= 0 {
			sendJSONError(w, "id error format", http.StatusBadRequest)
			return
		}

		task, err := db.GetTaskByID(dbase, id)
		if err != nil {
			sendJSONError(w, "task not found", http.StatusNotFound)
			return
		}

		if task.Repeat == "" {
			err = db.DeleteTask(dbase, id)
			if err != nil {
				sendJSONError(w, "Failed to delete task", http.StatusInternalServerError)
				return
			}
		} else {
			nextDate, err := logic.NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				sendJSONError(w, "Failed to calculate next date", http.StatusBadRequest)
				return
			}
			task.Date = nextDate

			err = db.UpdateTask(dbase, task)
			if err != nil {
				sendJSONError(w, "Failed to update task", http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(map[string]string{})
	}
}

// deleteTask handles the deletion of a task from the database by its ID.
// It reads the "id" parameter from the request, validates it, and removes the task if it exists.
// If the task is not found or the ID is invalid, it returns an appropriate JSON error response.
func deleteTask(dbase *sqlx.DB, w http.ResponseWriter, req *http.Request) {
	idStr := req.URL.Query().Get("id")
	if idStr == "" {
		sendJSONError(w, "id is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		sendJSONError(w, "id error format", http.StatusBadRequest)
		return
	}

	err = db.DeleteTask(dbase, id)
	if err != nil {
		sendJSONError(w, "Task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(map[string]string{})
}

// sendJSONError sends a structured JSON response with an error message.
func sendJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// SignInHandler authentication handler
func SignInHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, `{"error": "Invalid method"}`, http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Password string `json:"password"`
	}

	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		sendJSONError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if request.Password != os.Getenv("TODO_PASSWORD") {
		sendJSONError(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateToken()
	if err != nil {
		sendJSONError(w, "Token error", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
