package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"go_final_project/internal/constants"
	"go_final_project/internal/handlers/common"
	"go_final_project/internal/models"
)

type TaskService interface {
	GetTask(id int64) (*models.Task, error)
	GetTasksWithFilter(filterType int, filterValue string) ([]*models.Task, error)
	CreateTask(task *models.Task) (int64, error)
	UpdateTask(task *models.Task) error
	SetTaskDone(taskId int64) error
	DeleteTask(Id int64) error
}

type TaskHandler struct {
	svc TaskService
}

func NewTaskHandler(service TaskService) *TaskHandler {
	return &TaskHandler{svc: service}
}

func (h *TaskHandler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.handleGetTask(w, r)
		case http.MethodPost:
			h.handlePostTask(w, r)
		case http.MethodPut:
			h.handlePutTask(w, r)
		case http.MethodDelete:
			h.handleDeleteTask(w, r)
		default:
			http.Error(w, constants.ErrUnsupportedMethod, http.StatusMethodNotAllowed)
		}
	}
}

func (h *TaskHandler) handleGetTask(w http.ResponseWriter, r *http.Request) {
	id, err := common.GetIdFromQuery(r)
	if err != nil {
		common.RespondWithError(w, err)
		return
	}

	task, err := h.svc.GetTask(id)
	if err != nil {
		common.RespondWithError(w, err)
		return
	}

	response := models.GetTaskResponseDTO{
		Id:      strconv.FormatInt(task.Id, 10),
		Date:    task.Date,
		Title:   task.Title,
		Comment: task.Comment,
		Repeat:  task.Repeat,
	}

	common.Respond(w, response)
}

func (h *TaskHandler) handlePostTask(w http.ResponseWriter, r *http.Request) {
	var taskDTO models.Task
	err := json.NewDecoder(r.Body).Decode(&taskDTO)
	if err != nil {
		common.RespondWithError(w, constants.ErrInvalidJson)
		return
	}

	taskId, err := h.svc.CreateTask(&taskDTO)
	if err != nil {
		common.RespondWithError(w, err)
		return
	}

	response := models.Response{Id: taskId}
	common.Respond(w, response)
}

func (h *TaskHandler) handlePutTask(w http.ResponseWriter, r *http.Request) {
	var taskDTO models.TaskPutDTO
	err := json.NewDecoder(r.Body).Decode(&taskDTO)
	if err != nil {
		common.RespondWithError(w, constants.ErrInvalidJson)
		return
	}
	taskId, err := strconv.ParseInt(taskDTO.Id, 10, 64)
	if err != nil {
		common.RespondWithError(w, constants.ErrGetTaskId)
		return
	}

	taskRequest := &models.Task{
		Id:      taskId,
		Date:    taskDTO.Date,
		Title:   taskDTO.Title,
		Comment: taskDTO.Comment,
		Repeat:  taskDTO.Repeat,
	}
	err = h.svc.UpdateTask(taskRequest)
	if err != nil {
		common.RespondWithError(w, err)
		return
	}

	common.SetJsonHeader(w)
	w.Write([]byte("{}"))
}

func (h *TaskHandler) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := common.GetIdFromQuery(r)
	if err != nil {
		common.RespondWithError(w, err)
		return
	}

	err = h.svc.DeleteTask(id)
	if err != nil {
		common.RespondWithError(w, err)
		return
	}

	common.SetJsonHeader(w)
	w.Write([]byte("{}"))
}
