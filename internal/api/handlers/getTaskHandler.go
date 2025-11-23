package handlers

import (
	"go-task-manager-final_project/internal/api"
	"go-task-manager-final_project/internal/db"
	"net/http"
	"strings"
)

// Обработчик HTTP-запроса для получения задачи по ID.
// Параметры:
// w - объект для записи HTTP-ответа;
// r - HTTP-запрос с параметрами.
// Логика:
//  1. Извлекает параметр id из запроса.
//  2. Проверяет наличие ID.
//  3. Запрашивает задачу из БД по ID.
//  4. Возвращает результат (задачу или ошибку).
func (s *APIServer) getTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем параметр id из строки запроса
	id := r.URL.Query().Get("id")

	// Проверяем, что ID не пустой
	if strings.TrimSpace(id) == "" {
		api.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "id parameter is required",
		})
		return
	}

	// Проверяем формат ID (числовой)
	if !api.IsValidID(id) {
		api.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid id format: must be a integer number",
		})
		return
	}

	// Вызываем БД для получения задачи по ID
	task, err := db.GetTask(s.DB, id)
	if err != nil {
		// Различаем типы ошибок для более точной обратной связи
		if err.Error() == "task with id "+id+" not found" {
			api.WriteJSON(w, http.StatusNotFound, map[string]string{
				"error": "task not found",
			})
			return
		}
		// Логируем неожиданную ошибку БД
		api.WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to fetch task from database",
		})
		return
	}

	// Формируем успешный ответ с найденной задачей
	// Статус: HTTP 200 OK
	// Тело ответа: объект задачи в JSON-формате.
	api.WriteJSON(w, http.StatusOK, task)
}
