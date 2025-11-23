package handlers

import (
	"encoding/json"
	"fmt"
	"go-task-manager-final_project/internal/api"
	"go-task-manager-final_project/internal/db"
	"net/http"
	"strings"
)

// putTaskHandler обрабатывает HTTP-запрос на обновление задачи.
// Параметры:
// w - объект http.ResponseWriter для отправки ответа клиенту;
// r - объект *http.Request с данными входящего запроса.
// Логика:
// - проверяет заголовок Content-Type на соответствие application/json;
// - декодирует JSON из тела запроса в структуру db.Task;
// - валидирует обязательные поля (например, Title);
// - проверяет и корректирует дату задачи;
// - обновляет задачу в базе данных;
// - возвращает соответствующий HTTP-ответ (ошибка или успех).
func (s *APIServer) putTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем значение заголовка Content-Type из запроса
	contentType := r.Header.Get("Content-Type")
	// Проверяем, что Content-Type начинается с "application/json" (без учёта регистра)
	if !strings.HasPrefix(strings.ToLower(contentType), "application/json") {
		api.WriteJSON(w, http.StatusUnsupportedMediaType, map[string]string{
			"error": "content-Type must be application/json",
		})
		return
	}

	// Создаём переменную для хранения данных задачи
	var task db.Task
	// Декодируем JSON из тела запроса в структуру task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		r.Body.Close()
		api.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": fmt.Sprintf("invalid JSON payload: %v", err),
		})
		return
	}

	// Проверяем, что поле Title не пустое (обязательное поле)
	if strings.TrimSpace(task.Title) == "" {
		api.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "title cannot be empty or whitespace",
		})
		return
	}

	// Проверяем и корректируем дату задачи (вызов вспомогательной функции)
	if err := checkDate(&task); err != nil {
		api.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Обновляем задачу в базе данных через функцию UpdateTask из пакета db
	err := db.UpdateTask(s.DB, &task)
	if err != nil {
		api.WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("failed to update task: %v", err),
		})
		return
	}

	// Отправляем успешный ответ с ID задачи, ссылкой на ресурс и сообщением
	api.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"id":       task.ID,
		"location": fmt.Sprintf("/tasks/%s", task.ID),
		"message":  "Task update successfully",
	})
}
