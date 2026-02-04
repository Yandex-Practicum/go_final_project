package controller

import (
	"context"
	"errors"
	"net/http"

	"todo/internal/entity"

	"github.com/go-chi/render"
)

type usecase interface {
	// AddTask adds a new task to the repository.
	// It takes a context, representing the execution context, and a task, representing the task to be added.
	// It returns the ID of the newly added task or an error if any.
	AddTask(ctx context.Context, task *entity.Task) (string, error)

	// GetTasks retrieves a list of tasks from the repository.
	// It takes a context, representing the execution context, and a query string to filter the tasks.
	// It returns a slice of Task entities and an error, if any.
	GetTasks(ctx context.Context, query string) ([]*entity.Task, error)

	// GetNextDate calculates the next date based on the current date, a given date, and a repeat pattern.
	// It takes the current date in the format specified by the entity.DateFormat constant,
	// and returns the next date in the same format.
	// If any error occurs during the calculation, it returns an empty string and the error.
	GetNextDate(now, date, repeat string) (string, error)

	// DeleteTask deletes a task with the specified ID.
	// It first checks if the task exists by calling the GetTask method of the repository.
	// If the task exists, it then calls the DeleteTask method of the repository to delete the task.
	// If any error occurs during the process, it returns an error with additional context information.
	DeleteTask(ctx context.Context, id string) error

	// DoTask performs a task based on the given task ID.
	// If the task has a repeat schedule, it calculates the next date and updates the task.
	// If the task does not have a repeat schedule, it deletes the task from the repository.
	// Returns an error if any operation fails.
	DoTask(ctx context.Context, id string) error

	// UpdateTask updates the given task in the repository.
	// It validates the task's ID, title, and date, and sets the date to the current time if it is empty.
	// It also calculates the next date based on the repeat interval, if provided.
	// If the calculated date is in the past, it updates the task's date to the next date.
	// Finally, it calls the repository's UpdateTask method to update the task.
	// If any validation or repository operation fails, it returns an error.
	UpdateTask(ctx context.Context, task *entity.Task) error

	// GetTask retrieves a task by its ID.
	// It takes a context and an ID as input parameters and returns the corresponding task and any error encountered.
	GetTask(ctx context.Context, id string) (*entity.Task, error)

	// Login authenticates the user with the provided password.
	// It returns a string representing the user's authentication token and an error if any.
	Login(password string) (string, error)

	// ValidateToken validates the given token using JWT validation.
	// It returns an error if the token is invalid or if JWT validation fails.
	ValidateToken(token string) error

	// ShouldCheckToken checks if the password is set and
	// returns a boolean value indicating whether token should be checked.
	ShouldCheckToken() bool
}

// TaskController is a controller implementation for task-related operations.
type TaskController struct {
	uc usecase
}

// NewTaskController creates a new instance of TaskController.
// It takes a usecase, representing the usecase to be used, as the parameter.
// It returns the initialized TaskController instance.
func NewTaskController(uc usecase) *TaskController {
	return &TaskController{
		uc: uc,
	}
}

func (tc *TaskController) sendError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, entity.ErrUnmarshal) ||
		errors.Is(err, entity.ErrEmptyTitle) ||
		errors.Is(err, entity.ErrInvalidDate):
		w.WriteHeader(http.StatusBadRequest)
	case errors.Is(err, entity.ErrTaskNotFound):
		w.WriteHeader(http.StatusNotFound)
	case errors.Is(err, entity.ErrInvalidPassword) ||
		errors.Is(err, entity.ErrAuthDisabled) ||
		errors.Is(err, entity.ErrAuthenticationRequired):
		w.WriteHeader(http.StatusUnauthorized)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	render.JSON(w, r, map[string]string{"error": err.Error()})
}
