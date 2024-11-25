package handler

import (
	"encoding/json"
	"net/http"

	"go_final_project/internal/error"
)

func (h *Handler) GetTaskById(w http.ResponseWriter, r *http.Request) {
	//получаем id
	id := r.URL.Query().Get("id")
	if id == "" {
		error.JsonResponse(w, "Отсутсвует указанный id")
		return
	}

	//Получаем значения полей задачи по id
	t, err := h.repo.GetTaskByID(id)
	if err != nil {
		error.JsonResponse(w, "Задача не найдена")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err = json.NewEncoder(w).Encode(t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
