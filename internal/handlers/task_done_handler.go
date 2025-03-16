package handlers

import (
	"net/http"

	"go_final_project/internal/constants"
	"go_final_project/internal/handlers/common"
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
			http.Error(w, constants.ErrUnsupportedMethod, http.StatusMethodNotAllowed)
		}
	}
}

func (h *TaskDoneHandler) handlePostTaskDone(w http.ResponseWriter, r *http.Request) {
	id, err := common.GetIdFromQuery(r)
	if err != nil {
		common.RespondWithError(w, err)
		return
	}

	err = h.svc.SetTaskDone(id)
	if err != nil {
		common.RespondWithError(w, err)
		return
	}

	common.SetJsonHeader(w)
	w.Write([]byte("{}"))
}
