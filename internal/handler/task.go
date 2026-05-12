package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"go_final_project/internal/models"
	"go_final_project/pkg/api"
)

const TaskCount = 50

var (
	ErrTaskParsed = errors.New("fail decode json")
)

type ITaskService interface {
	Add(*models.Task) (string, error)
	List(int) ([]models.Task, error)
	Get(string) (*models.Task, error)
	Edit(*models.Task) error
	Done(string) error
	Delete(string) error
	NextDate(string, string, string) (string, error)
}

type TaskHandler struct {
	service ITaskService
	logs    *log.Logger
}

func NewTaskHandler(service ITaskService, logs *log.Logger) *TaskHandler {
	return &TaskHandler{
		service: service,
		logs:    logs,
	}
}

func (h *TaskHandler) Add() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var task models.Task

		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			api.Failed(w, ErrTaskParsed, http.StatusBadRequest)
			return
		}

		id, err := h.service.Add(&task)
		if err != nil {
			api.Failed(w, err, http.StatusBadRequest)
			return
		}

		api.Succefull(w, models.TaskCreated{
			ID: id,
		})
	})
}

func (h *TaskHandler) List() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tasks, err := h.service.List(TaskCount)
		if err != nil {
			api.Failed(w, err, http.StatusBadRequest)
			return
		}

		api.Succefull(w, models.Tasks{
			Tasks: tasks,
		})
	})
}

func (h *TaskHandler) Get() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		task, err := h.service.Get(r.FormValue("id"))
		if err != nil {
			api.Failed(w, err, http.StatusBadRequest)
			return
		}

		api.Succefull(w, task)
	})
}

func (h *TaskHandler) Edit() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var task models.Task

		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			api.Failed(w, ErrTaskParsed, http.StatusBadRequest)
			return
		}

		err = h.service.Edit(&task)
		if err != nil {
			api.Failed(w, err, http.StatusBadRequest)
			return
		}

		api.Succefull(w, struct{}{})
	})
}

func (h *TaskHandler) Done() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := h.service.Done(r.FormValue("id"))
		if err != nil {
			api.Failed(w, err, http.StatusBadRequest)
			return
		}

		api.Succefull(w, struct{}{})
	})
}

func (h *TaskHandler) Delete() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := h.service.Delete(r.FormValue("id"))
		if err != nil {
			api.Failed(w, err, http.StatusBadRequest)
			return
		}

		api.Succefull(w, struct{}{})
	})
}

func (h *TaskHandler) NextDate() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		date, err := h.service.NextDate(r.FormValue("now"), r.FormValue("date"), r.FormValue("repeat"))
		if err != nil {
			api.Failed(w, err, http.StatusBadRequest)
			return
		}

		_, err = io.WriteString(w, date)
		if err != nil {
			h.logs.Println(err)
		}
	})
}
