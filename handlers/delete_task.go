package handlers

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/FunnyFoXD/go_final_project/databases"
)

// DeleteTaskHandler is a handler for "/api/task" endpoint.
// It deletes a task with the given id from the database.
// It expects a GET request with a single parameter "id" which is a string.
// If the id is not given, it returns an error with HTTP status code 400.
// If the id is invalid, it returns an error with HTTP status code 400.
// If there is a database error, it returns an error with HTTP status code 500.
// If the task is successfully deleted, it returns an empty JSON object with HTTP status code 200.
func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	if id == "" {
		http.Error(w, `{"error":"identifier is empty"}`, http.StatusBadRequest)
		return
	}

	if !isValidID(id) {
		http.Error(w, `{"error":"invalid identifier"}`, http.StatusBadRequest)
		return
	}

	err := databases.DeleteTask(id)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte("{}"))
}

// isValidID checks whether the given string is a valid task identifier.
//
// A valid identifier is a non-empty string consisting only of digits.
func isValidID(id string) bool {
	matched, _ := regexp.MatchString(`^\d+$`, id)
	return matched
}
