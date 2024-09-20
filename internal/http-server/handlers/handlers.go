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
	AddTask(t tasks.Task) (int, error)
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

		err = task.Validate()
		if err != nil {
			log.Error("failed to validate task", logger.Err(err))
			render.JSON(w, r, api.NewErrResponse("Failed to validate request`s body data"))
			return
		}

		insertId, err := storage.AddTask(task)
		if err != nil {
			log.Error("failed to add task into storage", logger.Err(err))
			render.JSON(w, r, api.NewErrResponse("Failed to add task. Server internal error."))
		}

		render.JSON(w, r, api.NewPostSucceedResponse(insertId))
	}

}
