package handlers

import (
	"go_final_project/internal/models"
	"go_final_project/internal/utils"
	"net/http"
	"strconv"
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
			http.Error(w, utils.ErrUnsupportedMethod, http.StatusMethodNotAllowed)
		}
	}
}

func (h *TasksHandler) handleGetTasks(w http.ResponseWriter, r *http.Request) {
	filterType, filterValue := utils.GetFilterTypeAndValue(r)
	tasks, err := h.svc.GetTasksWithFilter(filterType, filterValue)
	if err != nil {
		utils.RespondWithError(w, err)
		return
	}

	response := models.GetTasksResponseDTO{Tasks: make([]models.GetTaskResponseDTO, 0)}
	for _, t := range tasks {
		response.Tasks = append(response.Tasks, models.GetTaskResponseDTO{
			ID:      strconv.FormatInt(t.ID, 10),
			Date:    t.Date,
			Title:   t.Title,
			Comment: t.Comment,
			Repeat:  t.Repeat,
		})
	}

	utils.Respond(w, response)
}
