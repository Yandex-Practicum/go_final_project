package handlers

import (
	"encoding/json"
	"final-project-bronner/go/models"
	"net/http"
)

type GetAllTasksResponse struct {
	Tasks []models.Task `json:"tasks"`
}

func (h *Handler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tasks, err := h.db.GetAllTasks()
	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	resultToSerialize := GetAllTasksResponse{Tasks: tasks}
	resp, err := json.Marshal(resultToSerialize)
	if err != nil {
		responseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
