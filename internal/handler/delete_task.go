package handler

import (
	"net/http"
	"strconv"

	"go_final_project/internal/error"
)

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {

	//Проверяем id
	id := r.URL.Query().Get("id")
	if id == "" {
		error.JsonResponse(w, "Отсутсвует указанный id")
		return
	}

	//Конвертируем id в int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		error.JsonResponse(w, "Ошибка конвертации id в Int")
		return
	}

	//Удаляем задачу по id
	err = h.repo.DeleteTask(idInt)
	if err != nil {
		error.JsonResponse(w, "Ошибка удаления задачи")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_, err = w.Write([]byte("{}"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
