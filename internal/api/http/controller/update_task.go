package controller

import (
	"net/http"

	"todo/internal/entity"

	"github.com/go-chi/render"
)

// UpdateTask updates a task based on the request payload.
// It decodes the request body into a task object,
// then calls the UpdateTask method of the TaskController's underlying use case.
// If there is an error during decoding or updating the task, it sends an appropriate error response.
// Finally, it responds with an empty JSON object.
func (tc *TaskController) UpdateTask(w http.ResponseWriter, r *http.Request) {
	task := &entity.Task{}
	if err := render.Decode(r, task); err != nil {
		tc.sendError(w, r, entity.ErrUnmarshal)
		return
	}

	if err := tc.uc.UpdateTask(r.Context(), task); err != nil {
		tc.sendError(w, r, err)
		return
	}

	render.JSON(w, r, map[string]any{})
}
