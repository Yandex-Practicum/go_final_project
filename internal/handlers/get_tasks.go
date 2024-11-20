package handlers

import (
	"encoding/json"
	"final_project/internal/common"
	"net/http"
	"time"
)

func (h *Handler) GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if param := r.FormValue("search"); param == "" {
		tasks, err := h.rep.GetAllTasks()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(common.Response{Error: err.Error()})
			return
		}
		if tasks.Tasks == nil {
			tasks.Tasks = []common.Task{}
		}
		json.NewEncoder(w).Encode(tasks)
		return
	} else if dateTime, err := time.Parse("02.01.2006", param); err == nil {
		tasks, err := h.rep.GetTasksByDate(dateTime.Format(common.TimeFormat))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(common.Response{Error: err.Error()})
			return
		}
		if tasks.Tasks == nil {
			tasks.Tasks = []common.Task{}
		}
		json.NewEncoder(w).Encode(tasks)
		return
	} else if param != "" {

		tasks, err := h.rep.GetTasksByParam(param)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(common.Response{Error: err.Error()})
			return
		}
		if tasks.Tasks == nil {
			tasks.Tasks = []common.Task{}
		}

		json.NewEncoder(w).Encode(tasks)
		return
	}
}
