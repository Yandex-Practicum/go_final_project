package handlers

import (
	"database/sql"
	"fmt"
	"go-task-manager-final_project/internal/api"
	"go-task-manager-final_project/internal/db"

	"net/http"
	"strings"
)

// Обработчик HTTP-запроса на удаление задачи.
// Параметры:
// w - объект для записи HTTP-ответа;
// r - HTTP-запрос с информацией о запросе (включая параметры URL).
// Логика:
//  1. Извлекает параметр id из строки запроса.
//  2. Проверяет, что id не пустой.
//  3. Пытается удалить задачу по указанному id.
//  4. Возвращает соответствующий HTTP-статус и JSON-ответ в зависимости от результата.
func (s *APIServer) deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем параметр id из строки запроса (например, /delete?id=123)
	id := r.URL.Query().Get("id")

	// Проверяем, что ID не пустой и не состоит только из пробелов
	if strings.TrimSpace(id) == "" {
		api.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "missing id parameter",
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

	// Пытаемся удалить задачу с указанным ID из базы данных
	err := db.DeleteTask(s.DB, id)
	if err != nil {
		// Если задача не найдена в БД (стандартная ошибка SQL), возвращаем статус 404 (Not Found)
		if err == sql.ErrNoRows {
			api.WriteJSON(w, http.StatusNotFound, map[string]string{
				"error": "task not found in database",
			})
		} else {
			// Любая другая ошибка при удалении (например, проблемы с соединением), возвращаем статус 500 (Internal Server Error)
			api.WriteJSON(w, http.StatusInternalServerError, map[string]string{
				"error": fmt.Sprintf("could not delete task: %v", err),
			})
		}
		return
	}

	// Если удаление прошло успешно - возвращаем пустой JSON-объект и статус 200 (OK)
	api.WriteJSON(w, http.StatusOK, map[string]interface{}{})
}
