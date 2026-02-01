package api

import (
	"github.com/go-chi/chi/v5"
)

func Init(r chi.Router) {
	r.Get("/api/nextdate", getNextDayHandler)
	r.With(auth).Post("/api/task", addTaskHandler)
	r.With(auth).Get("/api/tasks", getTaskListHandler)
	r.With(auth).Get("/api/task", getTaskByIdHandler)
	r.With(auth).Put("/api/task", updateTaskHandler)
	r.With(auth).Delete("/api/task", deleteTaskHandler)
	r.With(auth).Post("/api/task/done", completeTaskHandler)
	r.Post("/api/signin", checkUser)
}
