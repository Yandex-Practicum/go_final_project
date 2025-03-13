package controller

import (
	"net/http"

	"github.com/go-chi/render"
)

// GetTasks retrieves tasks based on the provided search query.
// It takes a http.ResponseWriter and a http.Request as parameters.
// The search query is extracted from the request URL query parameters.
// It calls the GetTasks method of the TaskController's underlying use case to fetch the tasks.
// If an error occurs during the retrieval process, it sends an error response.
// Otherwise, it sends a JSON response containing the retrieved tasks.
func (tc *TaskController) GetTasks(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("search")
	tasks, err := tc.uc.GetTasks(r.Context(), query)
	if err != nil {
		tc.sendError(w, r, err)
		return
	}

	render.JSON(w, r, map[string]any{"tasks": tasks})
}
