package handlers

import (
	"net/http"
	"strconv"

	"go_final_project/internal/constants"
	"go_final_project/internal/handlers/common"
	"go_final_project/internal/models"
)

type TasksHandler struct {
	svc TaskService
}

func NewTasksHandler(service TaskService) *TasksHandler {
	return &TasksHandler{svc: service}
}

func (h *TasksHandler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.handleGetTasks(w, r)
		default:
			http.Error(w, constants.ErrUnsupportedMethod, http.StatusMethodNotAllowed)
		}
	}
}

func (h *TasksHandler) handleGetTasks(w http.ResponseWriter, r *http.Request) {
	filterType, filterValue := common.GetFilterTypeAndValue(r)
	tasks, err := h.svc.GetTasksWithFilter(filterType, filterValue)
	if err != nil {
		common.RespondWithError(w, err)
		return
	}

	response := models.GetTasksResponseDTO{Tasks: make([]models.GetTaskResponseDTO, 0)}
	for _, t := range tasks {
		response.Tasks = append(response.Tasks, models.GetTaskResponseDTO{
			Id:      strconv.FormatInt(t.Id, 10),
			Date:    t.Date,
			Title:   t.Title,
			Comment: t.Comment,
			Repeat:  t.Repeat,
		})
	}

	common.Respond(w, response)
}
