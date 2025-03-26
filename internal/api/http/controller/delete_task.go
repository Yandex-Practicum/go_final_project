package controller

import (
	"net/http"

	"github.com/go-chi/render"
)

// DeleteTask deletes a task with the given ID.
// It takes a http.ResponseWriter and a http.Request as parameters.
// It parses the ID from the request URL query parameters and deletes the task using the TaskController's use case.
// If there is an error parsing the ID or deleting the task, it sends an error response.
// Finally, it sends an empty JSON response indicating successful deletion.
func (tc *TaskController) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if err := tc.uc.DeleteTask(r.Context(), id); err != nil {
		tc.sendError(w, r, err)
		return
	}

	render.JSON(w, r, map[string]any{})
}
