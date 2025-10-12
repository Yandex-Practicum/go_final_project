package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go_final_project/pkg/db"
)

// GetTaskHandler обработчик для получения задачи по ID
func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, map[string]string{"error": "Метод не поддерживается"}, http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор"}, http.StatusBadRequest)
		return
	}

	fmt.Printf("Получен запрос на получение задачи ID: %s\n", id)

	task, err := db.GetTask(id)
	if err != nil {
		// УБЕДИТЕСЬ ЧТО НЕТ ДВОЙНОГО WriteHeader
		writeJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	writeJSON(w, task, http.StatusOK)
}

// UpdateTaskHandler обработчик для обновления задач
func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeJSON(w, map[string]string{"error": "Метод не поддерживается"}, http.StatusMethodNotAllowed)
		return
	}

	var task db.Task

	// Декодируем JSON
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, map[string]string{"error": fmt.Sprintf("Ошибка декодирования JSON: %v", err)}, http.StatusBadRequest)
		return
	}

	fmt.Printf("получен запрос на обновление ID: %s\n", task.ID)

	// Валидируем обязаьтедбные поля
	if task.ID == "" {
		writeJSON(w, map[string]string{"error": "Не указан идинтификатор задачи"}, http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		writeJSON(w, map[string]string{"error": "Не указан заголовок задачи"}, http.StatusBadRequest)
		return
	}

	// Обрабатываем дату и правило повторения
	if err := processTaskDate(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	// Обновляем задачу в БД
	if err := db.UpdateTask(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	fmt.Printf("Задача ID %s успешно обновлена\n", task.ID)

	// Возвращаем пустой JSON при успехе
	writeJSON(w, map[string]interface{}{}, http.StatusOK)

}
