package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

// обработчик HTTP-запроса для получения информации об одной конкретной задаче по ее идентификатору
func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.URL.Query().Get("id") //Получение идентификатора задачи
	if id == "" {                 //Проверка наличия идентификатора
		http.Error(w, `{"error": "Не указан идентификатор"}`, http.StatusBadRequest)
		return
	}

	task, err := h.repo.GetTaskByID(id) //Поиск задачи в репозитории
	if err != nil {
		log.Println("Error getting task:", err)
		http.Error(w, `{"error": "Ошибка получения задачи"}`, http.StatusInternalServerError)
		return
	}

	if task == nil { //Проверка наличия найденной задачи
		http.Error(w, `{"error": "Задача не найдена"}`, http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(task) //Формирование JSON-ответа
	if err != nil {
		log.Println("Error encoding JSON:", err)
		http.Error(w, `{"error": "Ошибка кодирования JSON"}`, http.StatusInternalServerError)
		return
	}

}
