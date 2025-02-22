package handlers

import (
	"net/http"
	"strconv"

	"go_final/app/database"
)

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, `{"error":"Task id not specified"}`, http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, `{"error":"Incorrect task id"}`, http.StatusBadRequest)
		return
	}

	err = database.DeleteTaskByID(id)
	if err != nil {
		http.Error(w, `{"error":"Error deleting task"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{}`))
}
