package api

import (
	"net/http"
)

// Обработчик для /api/task, который определяет действие в зависимости от HTTP-метода
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r) // Добавление новой задачи
	case http.MethodGet:
		getTaskHandler(w, r) // Получение задачи по id
	case http.MethodPut:
		updateTaskHandler(w, r) // Обновление задачи
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}
