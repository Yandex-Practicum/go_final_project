package controller

import (
	"net/http"

	"github.com/go-chi/render"
)

// GetTask retrieves a task by its ID and sends it as a JSON response.
func (tc *TaskController) GetTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	task, err := tc.uc.GetTask(r.Context(), id)
	if err != nil {
		tc.sendError(w, r, err)
		return
	}

	render.JSON(w, r, task)
}
