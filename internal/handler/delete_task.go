package handler

import (
	"encoding/json"
	"net/http"
)

// Обработчик HTTP-запроса для удаления задачи по ее идентификатору
func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := r.URL.Query().Get("id") //Получение идентификатора задачи
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Не указан идентификатор задачи"})
		return
	}

	err := h.repo.DeleteTask(idStr) //Удаление задачи
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Ошибка удаления задачи: " + err.Error()})
		return
	}

	// Возвращаем пустой JSON
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{})
}
