package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/askoren1/go_final_project/internal/models"
)

// обработчик HTTP-запроса для получения списка всех задач
func (h *Handler) GetList(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	tasks, err := h.repo.GetList() //Вызываем метод GetList репозитория для получения списка всех задач из БД

	if err != nil {
		log.Println("Error getting tasks:", err)
		http.Error(w, `{"error": "Ошибка получения списка задач"}`, http.StatusInternalServerError)
		return
	}

	// подготовка ответа
	if tasks == nil {
		tasks = []models.Task{} //Проверка на пустой список
	}

	resp := struct {
		Tasks []models.Task `json:"tasks"` //Формирование структуры для ответа
	}{
		Tasks: tasks,
	}

	//Кодирование JSON-ответа
	err = json.NewEncoder(w).Encode(resp)

	if err != nil {
		log.Println("Error encoding JSON:", err)
		http.Error(w, `{"error": "Ошибка кодирования JSON"}`, http.StatusInternalServerError)
		return
	}

}
