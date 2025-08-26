package taskhandlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"go_final_project/internal/domain/entities"
	"go_final_project/internal/domain/services"
)

// AddTaskHandler - handler for /api/task
func AddTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlePostTask(w, r, db)
		case http.MethodPut:
			HandleUpdateTask(w, r, db)
		case http.MethodDelete:
			DeleteTaskHandler(db)(w, r)
		default:
			http.Error(w, `{"error":"Method not supported"}`, http.StatusMethodNotAllowed)
		}
	}
}

// handlePostTask - processes POST request (adding a task)
func handlePostTask(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	log.Println("New task addition request received")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var task entities.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		log.Printf("JSON parsing error: %v", err)
		http.Error(w, `{"error":"JSON parsing error"}`, http.StatusBadRequest)
		return
	}

	log.Printf("Task data received: %+v", task)

	// Check if task title is provided
	if task.Title == "" {
		log.Println("Error: Task title is missing")
		http.Error(w, `{"error":"Task title is required"}`, http.StatusBadRequest)
		return
	}

	// Get the current date in "YYYYMMDD" format
	nowStr := time.Now().Format("20060102")

	// If the date is not specified, use today’s date
	if task.Date == "" || task.Date == "today" {
		task.Date = nowStr
		log.Printf("Task date not specified, setting to today: %s", task.Date)
	}

	// Validate the date format
	taskDate, err := time.Parse("20060102", task.Date)
	if err != nil {
		log.Printf("Error: Invalid date format %s", task.Date)
		http.Error(w, `{"error":"Invalid date format"}`, http.StatusBadRequest)
		return
	}
	log.Printf("Initial task date: %s", taskDate.Format("2006-01-02"))

	// If the date is already today, do nothing
	if task.Date == nowStr {
		log.Printf("Task date is already today: %s, no changes made", task.Date)
	} else if task.Date < nowStr {
		// If the date is in the past, process repetition rules
		log.Printf("Task date is in the past (%s), applying processing...", task.Date)

		if task.Repeat == "" {
			// No repetition, set the date to today
			task.Date = nowStr
			log.Printf("No repetition specified, setting date to today: %s", task.Date)
		} else {
			// If repetition is specified, calculate the next available date
			nextDate, err := services.NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				log.Printf("Repetition error: %v", err)
				http.Error(w, fmt.Sprintf(`{"error":"Repetition error: %s"}`, err.Error()), http.StatusBadRequest)
				return
			}
			task.Date = nextDate
			log.Printf("Repetitive task, next occurrence: %s", task.Date)
		}
	} else {
		log.Printf("Task date is in the future: %s", task.Date)
	}

	// Validate repeat value
	if task.Repeat != "" && !isValidRepeat(task.Repeat) {
		log.Printf("Error: Invalid repeat value (%s)", task.Repeat)
		http.Error(w, `{"error":"Invalid repeat value"}`, http.StatusBadRequest)
		return
	}

	// Insert task into the database
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		log.Printf("Error adding task to database: %v", err)
		http.Error(w, `{"error":"Error adding task"}`, http.StatusInternalServerError)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("Error retrieving task ID: %v", err)
		http.Error(w, `{"error":"Error retrieving task ID"}`, http.StatusInternalServerError)
		return
	}

	log.Printf("Task added with ID: %d, date: %s", id, task.Date)
	json.NewEncoder(w).Encode(map[string]interface{}{"id": id})
}

// HandleUpdateTask - processes task updates
func HandleUpdateTask(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var task entities.Task

	// Parse JSON
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, `{"error":"JSON parsing error"}`, http.StatusBadRequest)
		return
	}

	// Check if ID is provided
	if task.ID == 0 {
		http.Error(w, `{"error":"Task ID is required"}`, http.StatusBadRequest)
		return
	}

	// Check if task exists
	var existingID int
	err := db.QueryRow("SELECT id FROM scheduler WHERE id = ?", task.ID).Scan(&existingID)
	if err == sql.ErrNoRows {
		http.Error(w, `{"error":"Task not found"}`, http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, `{"error":"Error checking task"}`, http.StatusInternalServerError)
		return
	}

	// Validate task details
	if task.Title == "" {
		http.Error(w, `{"error":"Task title cannot be empty"}`, http.StatusBadRequest)
		return
	}
	if _, err := time.Parse("20060102", task.Date); err != nil {
		http.Error(w, `{"error":"Invalid date format"}`, http.StatusBadRequest)
		return
	}
	if task.Repeat != "" && !isValidRepeat(task.Repeat) {
		http.Error(w, `{"error":"Invalid repeat value"}`, http.StatusBadRequest)
		return
	}

	// Update task in the database
	query := `UPDATE scheduler SET date=?, title=?, comment=?, repeat=? WHERE id=?`
	_, err = db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		http.Error(w, `{"error":"Error updating task"}`, http.StatusInternalServerError)
		return
	}

	log.Printf("Task with ID %d successfully updated", task.ID)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{}`))
}

// isValidRepeat - checks if the repeat value is valid
func isValidRepeat(repeat string) bool {
	allowedPrefixes := []string{"d ", "w ", "m ", "y"}
	for _, prefix := range allowedPrefixes {
		if strings.HasPrefix(repeat, prefix) {
			return true
		}
	}
	return false
}
