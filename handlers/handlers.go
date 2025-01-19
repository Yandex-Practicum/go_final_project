package handlers

import (
	"database/sql"
	"encoding/json"
	"go_final_project/models"
	"net/http"
	"strconv"
	"time"
)

const DefaultTaskLimit = 50

type Handler struct {
	DB *sql.DB
}

type TaskListResponse struct {
	Tasks []models.Task `json:"tasks"`
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{DB: db}
}

func (h *Handler) HandleTaskList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	limit := DefaultTaskLimit
	queryLimit := r.URL.Query().Get("limit")
	if queryLimit != "" {
		if parsedLimit, err := strconv.Atoi(queryLimit); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	rows, err := h.DB.Query(
		"SELECT id, date, title, comment, repeat, completed FROM scheduler ORDER BY date LIMIT ?",
		limit,
	)
	if err != nil {
		writeError(w, "Failed to retrieve tasks")
		return
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		var id int64
		var dateStr string
		err := rows.Scan(&id, &dateStr, &task.Title, &task.Comment, &task.Repeat, &task.Completed)
		if err != nil {
			writeError(w, "Failed to parse tasks")
			return
		}
		task.ID = strconv.FormatInt(id, 10)

		parsedDate, err := time.Parse(constants.DateFormat, dateStr)
		if err != nil {
			writeError(w, "Failed to parse date")
			return
		}
		task.Date = parsedDate

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		writeError(w, "Error iterating over rows")
		return
	}

	if tasks == nil {
		tasks = []models.Task{}
	}

	response := TaskListResponse{Tasks: tasks}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		writeError(w, "Failed to encode tasks")
	}
}

func (h *Handler) HandleAddTask(w http.ResponseWriter, r *http.Request) {
	// Implementation goes here
}

func (h *Handler) HandleGetTask(w http.ResponseWriter, r *http.Request) {
	// Implementation goes here
}

func (h *Handler) HandleUpdateTask(w http.ResponseWriter, r *http.Request) {
	// Implementation goes here
}

func (h *Handler) HandleDeleteTask(w http.ResponseWriter, r *http.Request) {
	// Implementation goes here
}

func (h *Handler) HandleMarkTaskDone(w http.ResponseWriter, r *http.Request) {
	// Implementation goes here
}

func writeError(w http.ResponseWriter, message string) {
	http.Error(w, message, http.StatusInternalServerError)
}
