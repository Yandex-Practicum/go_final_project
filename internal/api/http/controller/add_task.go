package controller

import (
	"net/http"

	"todo/internal/entity"

	"github.com/go-chi/render"
)

// AddTask adds a new task to the system.
// It decodes the JSON request body into a Task struct,
// then calls the AddTask method of the underlying use case to add the task.
// If there is an error during decoding or adding the task, it sends an appropriate error response.
// Finally, it responds with the task ID in JSON format.
func (tc *TaskController) AddTask(w http.ResponseWriter, r *http.Request) {
	var task entity.Task

	if err := render.DecodeJSON(r.Body, &task); err != nil {
		tc.sendError(w, r, entity.ErrUnmarshal)
		return
	}

	taskID, err := tc.uc.AddTask(r.Context(), &task)
	if err != nil {
		tc.sendError(w, r, err)
		return
	}

	render.JSON(w, r, map[string]string{"id": taskID})
}
