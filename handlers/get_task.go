package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/FunnyFoXD/go_final_project/databases"
)

// GetTaskHandler is a handler for "/api/task" endpoint.
// It expects a GET request with a parameter "id" which is a string.
// If the id is not given, it returns an error with HTTP status code 400.
// If the id is invalid, it returns an error with HTTP status code 400.
// If there is a database error, it returns an error with HTTP status code 500.
// If the task is successfully retrieved, it returns a JSON object with HTTP status code 200.
func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	if id == "" {
		http.Error(w, `{"error":"identifier is empty"}`, http.StatusBadRequest)
		return
	}

	task, err := databases.GetTaskByID(id)
	if err == sql.ErrNoRows {
		http.Error(w, `{"error":"task not found"}`, http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusInternalServerError)
	}
}