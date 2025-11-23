package handlers

import (
	"database/sql"
	"fmt"
	"go-task-manager-final_project/internal/api"
	"go-task-manager-final_project/internal/db"
	"go-task-manager-final_project/internal/services"
	"net/http"
	"strings"
	"time"
)

// doneTaskHandler обрабатывает запрос на завершение задачи.
// В зависимости от наличия правила повторения (task.Repeat) либо удаляет задачу, либо вычисляет и устанавливает новую дату выполнения.
// Параметры:
// w - http.ResponseWriter для отправки ответа клиенту;
// r - *http.Request, входящий HTTP-запрос.
func (s *APIServer) doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем параметр id из строки запроса
	id := r.URL.Query().Get("id")

	// Проверяем, что ID не пустой и не состоит только из пробелов
	if strings.TrimSpace(id) == "" {
		api.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "id parameter required",
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

	// Пытаемся получить задачу из базы данных по указанному ID
	task, err := db.GetTask(s.DB, id)
	if err != nil {
		if err == sql.ErrNoRows {
			// Задача с таким ID не найдена в БД - возвращаем 404 (Not Found)
			api.WriteJSON(w, http.StatusNotFound, map[string]string{
				"error": "task not found",
			})
		} else {
			// Произошла непредвиденная ошибка БД - возвращаем 500 (Internal Server Error)
			api.WriteJSON(w, http.StatusInternalServerError, map[string]string{
				"error": "could not retrieve task from database",
			})
		}
		return
	}

	// Проверяем наличие правила повторения задачи
	// Если Repeat пуст - задача не периодическая, её нужно удалить
	if task.Repeat == "" {
		// Пытаемся удалить задачу из БД
		err = db.DeleteTask(s.DB, id)
		if err != nil {
			if err == sql.ErrNoRows {
				// Задача уже удалена или не существует - возвращаем 404 (Not Found)
				api.WriteJSON(w, http.StatusNotFound, map[string]string{
					"error": "task not found",
				})
			} else {
				// Неожиданная ошибка при удалении - возвращаем 500 (Internal Server Error)
				api.WriteJSON(w, http.StatusInternalServerError, map[string]string{
					"error": "could not delete task",
				})
			}
			return
		}
		// Успешное удаление - возвращаем 200 (OK) с пустым JSON-объектом
		api.WriteJSON(w, http.StatusOK, map[string]interface{}{})
		return
	}

	// Задача периодическая - нужно вычислить следующую дату выполнения
	// Используем текущую дату, дату задачи и правило повторения
	next, err := services.NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		// Ошибка при расчёте даты (например, некорректный формат Repeat) - возвращаем 400
		api.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": fmt.Sprintf("invalid repeat pattern: %v", err),
		})
		return
	}

	// Обновляем дату задачи в БД на вычисленную следующую дату
	err = db.UpdateDate(s.DB, next, id)
	if err != nil {
		// Ошибка при обновлении даты в БД - возвращаем 500 (Internal Server Error)
		api.WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "could not update task date",
		})
		return
	}

	// Успешное обновление задачи - возвращаем OK с пустым JSON-объектом
	api.WriteJSON(w, http.StatusOK, map[string]interface{}{})
}
