package handlers

import (
	"go_final_project/internal/utils"
	"net/http"
)

type TaskDoneHandler struct {
	svc TaskService
}

func NewTaskDoneHandler(service TaskService) *TaskDoneHandler {
	return &TaskDoneHandler{svc: service}
}

func (h *TaskDoneHandler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.handlePostTaskDone(w, r)
		default:
			http.Error(w, utils.ErrUnsupportedMethod, http.StatusMethodNotAllowed)
		}
	}
}

func (h *TaskDoneHandler) handlePostTaskDone(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDFromQuery(r)
	if err != nil {
		utils.RespondWithError(w, err)
		return
	}

	err = h.svc.SetTaskDone(id)
	if err != nil {
		utils.RespondWithError(w, err)
		return
	}

	utils.SetJsonHeader(w)
	w.Write([]byte("{}"))
}
