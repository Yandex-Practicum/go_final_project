package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type controller interface {
	// AddTask adds a new task to the system.
	// It decodes the JSON request body into a Task struct,
	// then calls the AddTask method of the underlying use case to add the task.
	// If there is an error during decoding or adding the task, it sends an appropriate error response.
	// Finally, it responds with the task ID in JSON format.
	AddTask(w http.ResponseWriter, r *http.Request)

	// GetTasks retrieves tasks based on the provided search query.
	// It takes a http.ResponseWriter and a http.Request as parameters.
	// The search query is extracted from the request URL query parameters.
	// It calls the GetTasks method of the TaskController's underlying use case to fetch the tasks.
	// If an error occurs during the retrieval process, it sends an error response.
	// Otherwise, it sends a JSON response containing the retrieved tasks.
	GetTasks(w http.ResponseWriter, r *http.Request)

	// GetNextDate returns the next date based on the provided parameters.
	// It takes the current date and time, a specific date, and a repeat pattern as input.
	// If an error occurs during the process, it sends an error response.
	// Otherwise, it sends the next date as a plain text response.
	GetNextDate(w http.ResponseWriter, r *http.Request)

	// DoTask handles the HTTP request to perform a task.
	// It parses the task ID from the request URL query parameters,
	// calls the UseCase's DoTask method to perform the task,
	// and sends the response back to the client.
	DoTask(w http.ResponseWriter, r *http.Request)

	// DeleteTask deletes a task with the given ID.
	// It takes a http.ResponseWriter and a http.Request as parameters.
	// It parses the ID from the request URL query parameters and deletes the task using the TaskController's use case.
	// If there is an error parsing the ID or deleting the task, it sends an error response.
	// Finally, it sends an empty JSON response indicating successful deletion.
	DeleteTask(w http.ResponseWriter, r *http.Request)

	// UpdateTask updates a task based on the request payload.
	// It decodes the request body into a task object,
	// then calls the UpdateTask method of the TaskController's underlying use case.
	// If there is an error during decoding or updating the task, it sends an appropriate error response.
	// Finally, it responds with an empty JSON object.
	UpdateTask(w http.ResponseWriter, r *http.Request)

	// GetTask retrieves a task by its ID and sends it as a JSON response.
	GetTask(w http.ResponseWriter, r *http.Request)

	// Login handles the HTTP request for user login.
	// It parses the request body to get the password,
	// calls the Login method of the TaskController's use case,
	// and returns a JSON response with the generated token.
	Login(w http.ResponseWriter, r *http.Request)

	// ValidateTokenMiddleware is a middleware function that validates the token in the request cookie.
	// If the token is valid, it allows the request to proceed to the next handler.
	// If the token is invalid or missing, it sends an authentication required error response.
	ValidateTokenMiddleware(next http.HandlerFunc) http.HandlerFunc
}

type Router struct {
	mux *chi.Mux
}

func NewRouter(controller controller) *Router {
	router := chi.NewRouter()
	router.Handle("/*", http.FileServer(http.Dir("web")))

	apiRouter := chi.NewRouter()
	apiRouter.Use(render.SetContentType(render.ContentTypeJSON))

	apiRouter.Post("/signin", controller.Login)
	apiRouter.Post("/task", controller.ValidateTokenMiddleware(controller.AddTask))
	apiRouter.Get("/tasks", controller.ValidateTokenMiddleware(controller.GetTasks))
	apiRouter.Get("/nextdate", controller.GetNextDate)
	apiRouter.Post("/task/done", controller.ValidateTokenMiddleware(controller.DoTask))
	apiRouter.Delete("/task", controller.ValidateTokenMiddleware(controller.DeleteTask))
	apiRouter.Put("/task", controller.ValidateTokenMiddleware(controller.UpdateTask))
	apiRouter.Get("/task", controller.ValidateTokenMiddleware(controller.GetTask))

	router.Mount("/api", apiRouter)

	return &Router{
		mux: router,
	}
}

func (r *Router) Handler() http.Handler {
	return r.mux
}
