package controller

import (
	"net/http"

	"github.com/go-chi/render"
)

// DoTask handles the HTTP request to perform a task.
// It parses the task ID from the request URL query parameters,
// calls the UseCase's DoTask method to perform the task,
// and sends the response back to the client.
func (tc *TaskController) DoTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if err := tc.uc.DoTask(r.Context(), id); err != nil {
		tc.sendError(w, r, err)
		return
	}

	render.JSON(w, r, map[string]any{})
}
