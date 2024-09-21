package handlers

import (
	"log/slog"
	"net/http"
	"time"
	"todo-list/internal/http-server/api"
	"todo-list/internal/lib/logger"
	"todo-list/internal/tasks"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type taskStorage interface {
	AddTask(t *tasks.Task) (int, error)
	GetTasks() ([]tasks.Task, error)
	GetTask(taskId string) (*tasks.Task, error)
	UpdateTask(t *tasks.Task) error
}

func GetNextDate(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		log := log.With("request_id", middleware.GetReqID(r.Context()))

		paramNow := r.FormValue("now")
		if paramNow == "" {
			log.Error("failed to get the next date: parameter \"now\" is empty")
			http.Error(w, "Cant get the next date because parameter \"now\" is not defined.", http.StatusBadRequest)
			return
		}

		now, err := time.Parse("20060102", paramNow)
		if err != nil {
			log.Error("failed to parse time.Time from parameter \"now\" value", logger.Err(err))
			http.Error(w, "wrong value of \"now\" parameter", http.StatusBadRequest)
			return
		}

		date := r.FormValue("date")
		if date == "" {
			log.Error("failed to get the next date: parameter \"date\" is empty")
			http.Error(w, "Cant get the next date because parameter \"date\" is not defined.", http.StatusBadRequest)
			return
		}

		repeat := r.FormValue("repeat")
		if repeat == "" {
			log.Error("failed to get the next date: parameter \"repeat\" is empty")
			http.Error(w, "Cant get the next date because parameter \"repeat\" is not defined.", http.StatusBadRequest)
			return
		}

		newDate, err := tasks.NextDate(now, date, repeat)
		if err != nil {
			log.Error("failed to get the next date", logger.Err(err), slog.Attr{Key: "func", Value: slog.StringValue("api.GetNextDate")})
			http.Error(w, "failed to get the next date", http.StatusBadRequest)
			return
		}

		w.Write([]byte(newDate))

	}
}

func PostTask(log *slog.Logger, storage taskStorage) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		log := log.With(slog.Attr{Key: "request_id", Value: slog.StringValue(middleware.GetReqID(r.Context()))})

		var task tasks.Task
		err := render.DecodeJSON(r.Body, &task)
		if err != nil {
			log.Error("failed to decode r.Body", logger.Err(err))
			render.JSON(w, r, api.NewErrResponse("Cannot read request`s body"))
			return
		}

		err = task.Validate(false)
		if err != nil {
			log.Error("failed to validate task", logger.Err(err))
			render.JSON(w, r, api.NewErrResponse("Failed to validate request`s body data"))
			return
		}

		insertId, err := storage.AddTask(&task)
		if err != nil {
			log.Error("failed to add task into storage", logger.Err(err))
			render.JSON(w, r, api.NewErrResponse("Failed to add task. Server internal error."))
		}

		render.JSON(w, r, api.NewPostSucceedResponse(insertId))
	}
}

func GetTasks(log *slog.Logger, storage taskStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		log := log.With(slog.Attr{Key: "request_id", Value: slog.StringValue(middleware.GetReqID(r.Context()))})

		result := api.TasksResponse{}
		tasksList, err := storage.GetTasks()
		if err != nil {
			log.Error("failed to get tasks list from storage", logger.Err(err))
			render.JSON(w, r, api.NewErrResponse("Failed to get tasks list."))
			return
		}

		result.TaskList = tasksList
		render.JSON(w, r, result)
	}
}

func GetTask(log *slog.Logger, storage taskStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(slog.Attr{Key: "request_id", Value: slog.StringValue(middleware.GetReqID(r.Context()))})

		paramId := r.FormValue("id")
		if paramId == "" {
			log.Error("Id parameter is empty. Cannot get the task from db.")
			render.JSON(w, r, api.NewErrResponse("Cannot get the task data: id is empty"))
			return
		}

		task, err := storage.GetTask(paramId)
		if err != nil {
			log.Error("failed to get task from storage", logger.Err(err))
			render.JSON(w, r, api.NewErrResponse("Failed to get task from database: internal server error"))
			return
		}

		render.JSON(w, r, task)

	}
}

func PutTask(log *slog.Logger, storage taskStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(slog.Attr{Key: "request_id", Value: slog.StringValue(middleware.GetReqID(r.Context()))})

		task := tasks.Task{}
		err := render.DecodeJSON(r.Body, &task)
		if err != nil {
			log.Error("failed to decode requests body", logger.Err(err))
			render.JSON(w, r, api.NewErrResponse("Failed to read request`s body."))
			return
		}

		err = task.Validate(true)
		if err != nil {
			log.Error("task info is invalid", logger.Err(err))
			render.JSON(w, r, api.NewErrResponse("Task info is invalid."))
			return
		}

		err = storage.UpdateTask(&task)
		if err != nil {
			log.Error("failed to update task in storage", logger.Err(err))
			render.JSON(w, r, api.NewErrResponse("Failed to update the task. Internal server error."))
			return
		}

		render.JSON(w, r, struct{}{})

	}
}
