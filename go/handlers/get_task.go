package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	var id int
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}
	task, err := h.db.GetTask(id)
	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}
	serializedTask, err := json.Marshal(task)
	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(serializedTask)
}
