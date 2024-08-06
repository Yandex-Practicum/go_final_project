package handlers

import (
	"database/sql"
	"fmt"
	"go_final_project/tasks"
	"net/http"
	"strconv"
	"time"
)

// DoneTaskHandler takes request and mark task as done by ID
func DoneTaskHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		SendErrorResponse(w, "DoneTaskHandler: Method not allowed", http.StatusBadRequest)
		return
	}

	idTask := r.FormValue("id")
	if idTask == "" {
		SendErrorResponse(w, "DoneTaskHandler(): No ID provided", http.StatusBadRequest)
		return
	}

	idTaskParsed, _ := strconv.Atoi(idTask)

	var task tasks.Task
	var id int64
	err := db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", idTaskParsed).Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	task.Id = fmt.Sprint(id)
	if err == sql.ErrNoRows {
		SendErrorResponse(w, "DoneTaskHandler: Task not found", http.StatusNotFound)
		return
	} else if err != nil {
		SendErrorResponse(w, "DoneTaskHandler: Error retrieving task data", http.StatusInternalServerError)
		return
	}

	now := time.Now()
	// get task repetition date
	if task.Repeat != "" {
		newTaskDate, err := tasks.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			SendErrorResponse(w, "DoneTaskHandler: Invalid repeat pattern", http.StatusBadRequest)
			return
		}

		task.Date = newTaskDate
		// update task if repeat rule set

		query := "UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat =? WHERE id = ?"
		_, err = db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.Id)

		if err != nil {
			SendErrorResponse(w, "DoneTaskHandler: Task not found", http.StatusInternalServerError)
			return
		}
	} else {
		// delete task if repeat rule not set
		query := "DELETE FROM scheduler WHERE id = ?"
		task.Id = fmt.Sprint(id)
		result, err := db.Exec(query, task.Id)
		if err != nil {
			SendErrorResponse(w, "DoneTaskHandler: Error deleting task", http.StatusInternalServerError)
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			SendErrorResponse(w, "DoneTaskHandler: Unable to determine the number of rows affected after deleting a task", http.StatusInternalServerError)
			return
		} else if rowsAffected == 0 {
			SendErrorResponse(w, "DoneTaskHandler: Task not found", http.StatusInternalServerError)
			return
		}
	}
	// send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{}`))
}
